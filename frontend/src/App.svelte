<script lang="ts">
  const API_URL = import.meta.env.DEV
    ? "http://localhost:5689"
    : window.location.origin;

  let downloadSpeed = 0;
  let uploadSpeed = 0;
  let ping = 0;
  let pingTries = 25;
  let downloadSize = 100;
  let downloadChunks = 10;
  let uploadSize = 100;
  let uploadChunks = 10;

  let stop = false;

  const bytesToHumanReadable = (bytes: number) => {
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    if (bytes === 0) return "0 B";
    const i = parseInt(String(Math.floor(Math.log(bytes) / Math.log(1024))));
    return (bytes / Math.pow(1024, i)).toFixed(2) + " " + sizes[i];
  };

  const testDownloadSpeed = async (size: number, chunks: number) => {
    stop = false;

    const chunkSize = Math.floor(size / chunks);
    let remainingSize = size;
    const speeds: number[] = [];

    while (remainingSize > 0 && !stop) {
      const start = performance.now();

      const response = await fetch(`${API_URL}/download?size=${chunkSize}`);

      if (!response.ok) {
        throw new Error("Download failed");
      }

      const now = performance.now();
      const duration = now - start;
      const speed = (chunkSize / duration) * 1000;
      speeds.push(speed);

      speeds.sort((a, b) => a - b);
      const medianSpeed = speeds[Math.floor(speeds.length / 2)];
      downloadSpeed = medianSpeed;

      remainingSize -= chunkSize;
    }
  };

  const testUploadSpeed = async (size: number, chunks: number) => {
    stop = false;

    const chunkSize = Math.floor(size / chunks);
    const chunk = new Uint8Array(chunkSize);
    let remainingSize = size;
    const speeds: number[] = [];

    while (remainingSize > 0 && !stop) {
      const start = performance.now();

      const response = await fetch(`${API_URL}/upload?size=${chunkSize}`, {
        method: "POST",
        body: chunk,
      });

      if (!response.ok) {
        throw new Error("Upload failed");
      }

      const now = performance.now();
      const duration = now - start;
      const speed = (chunkSize / duration) * 1000;
      speeds.push(speed);

      speeds.sort((a, b) => a - b);
      const medianSpeed = speeds[Math.floor(speeds.length / 2)];
      uploadSpeed = medianSpeed;

      remainingSize -= chunkSize;
    }
  };

  const testPing = async (tries: number) => {
    const pings: number[] = [];
    const observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      const lastEntry = entries[entries.length - 1];
      pings.push(lastEntry.duration);

      pings.sort((a, b) => a - b);
      const medianPing = pings[Math.floor(pings.length / 2)];
      ping = medianPing;
    });
    observer.observe({ type: "resource", buffered: true });

    for (let i = 0; i < tries; i++) {
      await fetch(`${API_URL}/ping`);
    }

    observer.disconnect();
  };
</script>

<h1>Speed Test</h1>
<button on:click={() => (stop = true)}>Stop</button>
<br />
<br />
<button
  on:click={() => testDownloadSpeed(downloadSize * 1024 * 1024, downloadChunks)}
  >Test Download Speed</button
>
<input type="number" placeholder="Size in MB" bind:value={downloadSize} />
<input type="number" placeholder="Chunks" bind:value={downloadChunks} />
<p>Download Speed: {bytesToHumanReadable(downloadSpeed)}/s</p>

<button on:click={() => testUploadSpeed(uploadSize * 1024 * 1024, uploadChunks)}
  >Test Upload Speed</button
>
<input type="number" placeholder="Size in MB" bind:value={uploadSize} />
<input type="number" placeholder="Chunks" bind:value={uploadChunks} />
<p>Upload Speed: {bytesToHumanReadable(uploadSpeed)}/s</p>

<button on:click={() => testPing(pingTries)}>Test Ping</button>
<input type="number" placeholder="Tries" bind:value={pingTries} />
<p>Ping: {ping.toFixed(2)}ms</p>

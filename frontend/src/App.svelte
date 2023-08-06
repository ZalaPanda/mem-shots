<script lang="ts">
  let result = "?";

  const onApiTest = async (event: MouseEvent & { currentTarget: HTMLButtonElement }) => {
    const { currentTarget } = event;
    try {
      currentTarget.disabled = true;
      const response = await fetch("/api");
      result = await response.text();
    } catch (error) {
      result = String(error);
    } finally {
      currentTarget.disabled = false;
    }
  };
</script>

<main>
  <button on:click={onApiTest}>API test</button>
  <pre>{result}</pre>
</main>

<style lang="less">
  main {
    padding: 20px;
  }
  pre {
    padding: 10px;
  }
</style>

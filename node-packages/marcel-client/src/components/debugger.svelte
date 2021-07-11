<script>
  import machine from "../machine";

  let show = false;
  if (IS_DEV_MODE) show = true;

  window.toggleMarcelDebugger = () => (show = !show);
</script>

<style>
  .debugger {
    display: flex;
    flex-direction: column;
    position: absolute;
    top: 0;
    right: 0;
    height: 100vh;
    width: 40em;
  }

  pre {
    font-size: 10px;
    height: 100%;
  }

  button {
    width: 100%;
    border: 2px black solid;
    padding: 0.5em;
  }
</style>

{#if show}
  <div class="debugger">
    <button on:click={window.toggleMarcelDebugger}>Hide debugguer</button>
    <pre>
      State: {$machine.state} {'\n\n'}Context: {JSON.stringify($machine.context, 0, 2)}
    </pre>
  </div>
{/if}

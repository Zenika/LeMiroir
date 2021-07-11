<script>
  import { onMount } from "svelte";
  import decorate from "./iframe";
  export let plugin, stylesvar, uri;

  let prevProps, iframeElement, iframe;

  $: props = stylesvar && plugin && { ...plugin.frontend.props, stylesvar };

  $: iframeElement &&
    iframeElement.addEventListener("load", () => {
      iframe = decorate(iframeElement);
    });

  const sendProps = newProps => {
    iframe.postMessage({
      type: "propsChange",
      payload: { newProps, prevProps }
    });
    prevProps = newProps;
  };

  $: if (iframe) {
    iframe.addMessageListener(message => {
      console.debug(`Message from ${plugin.instanceId}`, message);
      if (message.type === "loaded") sendProps(props);
    });
  }

  $: if (props && iframe) sendProps(props);
</script>

<style>
  iframe {
    width: 100%;
    height: 100%;
    border: none;
  }
</style>

<svelte:options tag="marcel-widget" />

{#if plugin && uri}
  <iframe
    bind:this={iframeElement}
    title="{plugin.eltName}#{plugin.instanceId}"
    src={uri} />
{/if}

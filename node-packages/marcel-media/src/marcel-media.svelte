<script>
  import "marcel-widget";

  export let media, uri;

  $: col = media && 100 / media.cols;
  $: row = media && 100 / media.rows;
  $: pluginsURI = `${uri}/plugins`;
</script>

<style>
  :host,
  .media {
    position: relative;
    width: 100%;
    height: 100%;
  }

  * {
    box-sizing: border-box;
  }

  .widget {
    position: absolute;
    transition: width 1s ease-in-out, height 1s ease-in-out, left 1s ease-in-out,
      top 1s ease-in-out;
  }
</style>

<svelte:options tag="marcel-media" />

{#if media}
  <div
    class="media"
    style={`
      background-color: ${media.stylesvar['background-color']};
      color: ${media.stylesvar['primary-color']};
      font-family: ${media.stylesvar['font-family']};
    `}>
    {#each media.plugins as plugin}
      <marcel-widget
        style={`
          left: ${plugin.frontend.x * col}%;
          top: ${plugin.frontend.y * row}%;
          width: ${plugin.frontend.cols * col}%;
          height: ${plugin.frontend.rows * row}%;
        `}
        class="widget"
        {plugin}
        stylesvar={media.stylesvar}
        uri="{pluginsURI}/{plugin.eltName}/{plugin.instanceId}/" />
    {/each}
  </div>
{/if}

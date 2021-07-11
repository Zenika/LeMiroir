<script>
  import "bulma/css/bulma.css";
  import "marcel-media";
  import Debugger from "./components/debugger.svelte";
  import { Message } from "svelma";
  import machine, { maxLoadingStep } from "./machine";
  import AuthForm from "./components/auth-form.svelte";
  import Progress from './components/progress.svelte';
</script>

<style>
  .container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
  }
</style>

{#if $machine.state === 'mediaLoaded'}

  <marcel-media
    media={$machine.context.media}
    uri="{$machine.context.config.apiURI}medias/{$machine.context.media.id}" />

{:else}

  <div class="container">
    {#if $machine.state === 'loggedOut'}
      <AuthForm />
    {:else if $machine.state === 'noMedia'}
      <Message active title="{$machine.context.client.name} is not associated" type="is-info">
        This client ({$machine.context.client.name}) is not associated to any media.<br>
        You can associate the client in the backoffice.<br>
      </Message>
    {:else if $machine.state === 'fatalError'}
      <Message active title="Error!" type="is-danger">{$machine.context.error}</Message>
    {:else}
      <Progress type="is-info" value={$machine.context.loadingStep} max={maxLoadingStep} />
    {/if}
  </div>

{/if}

<Debugger />

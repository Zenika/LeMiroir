<script>
  import { Input, Button, Message } from 'svelma';
  import machine from "../machine";
  let login, password;

  const handleLogin = () => {
    machine.send({ type: "login", login, password })
  };
</script>

<style>
 form {
   display: grid;
   gap: 1em;
   align-items: center;
 }
</style>

{#if $machine.context.error}
  <Message active title="Error!" type="is-danger">{$machine.context.error}</Message>
{/if}

<form on:submit|preventDefault={handleLogin}>
  <Input bind:value={login} />
  <Input type="password" bind:value={password} />
  <Button nativeType="submit" disabled={$machine.state !== 'loggedOut'} type="is-primary">
    Login
  </Button>
</form>

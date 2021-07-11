<script>
  import { afterUpdate } from 'svelte'
  import { tweened } from 'svelte/motion'
  import { Progress } from "svelma"

  export let value, max, type, duration = 400

  let progressBar, container
  let loadingProgression = tweened(value, duration)

  afterUpdate(() => {
    if(container) progressBar = container.querySelector('progress')
  })

  $: if(progressBar) loadingProgression.set(value)

  $: if(progressBar) progressBar.value = $loadingProgression || 0
</script>

<div bind:this={container}>
  <Progress {max} value={loadingProgression} {type}/>
</div>

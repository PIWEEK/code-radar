<script>
 import { afterUpdate } from "svelte";
 export let file;
 let prevFile;

 let items = 10;

 function onClick() {
   if (items + 10 > file.history.length) {
     items = file.history.length;
   } else {
     items += 10;
   }
 }

 $: {
   if (file != prevFile){
     items = 10;
     prevFile = file;
   }
 }
</script>

<ul class="file-history uk-list uk-list-striped">
  {#each file.history.slice(0, items) as history, i}
    <li>
      [{history.date}] {history.user} +{history.added} -{history.deleted}
    </li>
  {/each}

</ul>
{#if items < file.history.length}
  <a href="#"
     on:click={onClick}
     
  >Load more...</a>
{/if}

<style>
 .file-history {
   font-size: 0.8rem;
 }
</style>

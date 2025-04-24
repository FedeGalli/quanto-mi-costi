<!-- Tile.svelte -->
<script lang="ts">
    import NumberFlow from "@number-flow/svelte";
    import { slide } from "svelte/transition";
    export let data;
    let total: number = 0;
    $: {
        total = 0;
        data.forEach((item: any) => {
            total += item.value;
        });
    }
</script>

<div class="p-1 rounded-lg shadow border text-center">
    <div class="flex justify-between font-bold text-xl border-b pb-1 mb-1">
        <span>{data[0].category}</span>
        <NumberFlow value={total} />
    </div>

    {#each data as item}
        <div
            class="flex justify-between text-sm"
            transition:slide={{ duration: 500 }}
        >
            <span>{item.name}</span>
            <NumberFlow value={item.value} />
        </div>
    {/each}
</div>

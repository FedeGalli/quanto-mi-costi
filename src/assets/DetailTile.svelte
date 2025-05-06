<!-- Tile.svelte -->
<script lang="ts">
    import NumberFlow from "@number-flow/svelte";
    import { slide } from "svelte/transition";
    export let data;
    export let title;
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
        <span>{title}</span>
        <NumberFlow
            value={total}
            format={{
                style: "currency",
                currency: "EUR",
                maximumFractionDigits: 0,
            }}
            locales={"it-IT"}
            class="text "
        />
    </div>

    {#each data as item}
        <div
            class="flex justify-between text-sm p-1"
            transition:slide={{ duration: 500 }}
        >
            <span>{item.name}</span>
            <NumberFlow
                value={item.value}
                format={{
                    style: "currency",
                    currency: "EUR",
                    maximumFractionDigits: 0,
                }}
                locales={"it-IT"}
                class="text"
            />
        </div>
    {/each}
</div>

<script lang="ts">
    import { collapsedDetailTiles } from "../store/collapsedDetailTiles";
    import NumberFlow from "@number-flow/svelte";
    import { slide } from "svelte/transition";

    export let data;
    export let element: any;
    export let title;

    let total: number = 0;
    $: {
        total = 0;
        data[element].forEach((item: any) => {
            total += item.value;
        });
    }

    function toggle() {
        collapsedDetailTiles.update((tiles: any) => {
            const updated = Object.fromEntries(
                Object.keys(tiles).map((key) => [key, false]),
            );
            updated[element] = true;
            return updated;
        });
    }

    $: isOpen = $collapsedDetailTiles[element];
</script>

<div class="w-full mb-3">
    <div
        class="bg-neutral-900 border-purple-400 p-3 rounded-lg shadow border text-center"
    >
        <div class="flex justify-between items-center font-bold text-xl">
            <button
                on:click={toggle}
                class="flex justify-between items-center w-full text-left"
            >
                <div class="flex items-center gap-2">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 320 512"
                        class="w-4 h-4 transition-transform duration-300"
                        class:rotate-90={isOpen}
                    >
                        <path
                            fill="currentColor"
                            d="M310.6 233.4c12.5 12.5 12.5 32.8 0 45.3l-192 192c-12.5 12.5-32.8 12.5-45.3 0s-12.5-32.8 0-45.3L242.7 256 73.4 86.6c-12.5-12.5-12.5-32.8 0-45.3s32.8-12.5 45.3 0l192 192z"
                        />
                    </svg>
                    <span>{title}</span>
                </div>
                <NumberFlow
                    value={total}
                    format={{
                        style: "currency",
                        currency: "EUR",
                        maximumFractionDigits: 0,
                    }}
                    locales={"it-IT"}
                    class="text-xl font-bold mb-1"
                />
            </button>
        </div>

        {#if isOpen}
            <div transition:slide={{ duration: 500 }}>
                <div class="border-b mt-1 mb-1"></div>
                {#each data[element] as item}
                    <div
                        class="flex justify-between text-sm p-0.5"
                        transition:slide={{ duration: 500 }}
                    >
                        <span class="font-bold">{item.name}</span>
                        <NumberFlow
                            value={item.value}
                            format={{
                                style: "currency",
                                currency: "EUR",
                                maximumFractionDigits: 0,
                            }}
                            locales={"it-IT"}
                        />
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

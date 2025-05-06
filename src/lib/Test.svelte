<script lang="ts">
    let finalTotal = 314822;
    let delta = 14822;
    let costs = [
        { label: "Taxes", value: 11146, icon: "🏠", delta: -4.7 },
        { label: "Notary", value: 1800, icon: "📜", delta: null },
        { label: "Bank", value: 3054, icon: "🏛️", delta: -0.95 },
    ];
</script>

<div class="bg-zinc-900 text-white rounded-2xl p-8 max-w-6xl mx-auto">
    <h2 class="text-3xl font-bold mb-6">Summary</h2>
    <div class="flex flex-wrap gap-8 justify-center w-full">
        <!-- Chart Placeholder -->
        <div
            class="flex justify-center items-center bg-zinc-800 rounded-2xl w-full max-w-sm h-64 flex-grow"
        >
            <div class="text-zinc-400 font-bold text-lg">Bace Cost</div>
        </div>

        <!-- Total Price Box -->
        <div
            class="bg-zinc-950 border border-zinc-700 p-6 rounded-xl text-center w-full max-w-sm flex-grow"
        >
            <div class="text-lg font-semibold mb-2">
                💰 <span class="font-bold">Final Total Price</span>
            </div>
            <div class="text-5xl font-extrabold">
                {finalTotal.toLocaleString()} €
            </div>
            <div class="text-green-400 text-sm mt-1">
                + {delta.toLocaleString()} € vs base price
            </div>
        </div>
    </div>

    <!-- Category Tiles -->
    <div class="flex gap-4 flex-wrap justify-between mt-8 w-full">
        {#each costs as cost}
            <div
                class="bg-zinc-950 border border-zinc-700 rounded-xl p-6 flex-1 text-center min-w-[150px]"
            >
                <div class="text-base mb-2">{cost.icon} {cost.label}</div>
                <div class="text-2xl font-bold">
                    {cost.value.toLocaleString()} €
                </div>
                {#if cost.delta !== null}
                    <div class="text-red-400 text-sm mt-1">
                        – {Math.abs(cost.delta).toFixed(2)}%
                    </div>
                {/if}
            </div>
        {/each}
    </div>
</div>

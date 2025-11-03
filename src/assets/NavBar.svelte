<script lang="ts">
    import { slide } from "svelte/transition";
    export let selectedTab = "summary";
    export let pro = false;
    export let getProFunction: () => void;
    import { _ } from "svelte-i18n";
</script>

<div
    class="flex justify-between mb-8 rounded-2xl p-2 gap-2"
    style="background: rgba(23, 23, 23, 1);"
>
    <button
        class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl text-gray-300 transition-all duration-300 ease-in-out hover:bg-gray-800 {selectedTab ===
        'summary'
            ? 'bg-purple-400/70 text-white'
            : ''}"
        on:click={() => (selectedTab = "summary")}
        type="button"
    >
        📊
        <span
            ><h1 class="text-xs font-bold leading-tight">
                {$_("navBar.recap")}
            </h1></span
        >
    </button>

    <button
        class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl text-gray-300 transition-all duration-300 ease-in-out hover:bg-gray-800 {selectedTab ===
        'base'
            ? 'bg-purple-400/70 text-white'
            : ''}"
        on:click={() => (selectedTab = "base")}
        type="button"
    >
        📦
        <span
            ><h1 class="text-xs font-bold leading-tight">
                {$_("navBar.costsDetail")}
            </h1></span
        >
    </button>

    <button
        class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl text-gray-300 transition-all duration-300 ease-in-out hover:bg-gray-800 {selectedTab ===
        'mortgage'
            ? 'bg-purple-400/70 text-white'
            : ''}"
        on:click={() => (selectedTab = "mortgage")}
        type="button"
        transition:slide={{ duration: 500 }}
    >
        💰
        <span
            ><h1 class="text-xs font-bold leading-tight">
                {$_("navBar.mortgage")}
            </h1></span
        >
    </button>

    {#if pro}
        <button
            class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl text-gray-300 transition-all duration-300 ease-in-out hover:bg-gray-800 {selectedTab ===
            'mortgage_compare'
                ? 'bg-purple-400/70 text-white'
                : ''}"
            on:click={() => (selectedTab = "mortgage_compare")}
            type="button"
            transition:slide={{ duration: 500 }}
        >
            ⚖️
            <span
                ><h1 class="text-xs font-bold leading-tight">
                    {$_("navBar.durationCompare")}
                </h1></span
            >
        </button>

        <button
            class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl text-gray-300 transition-all duration-300 ease-in-out hover:bg-gray-800 {selectedTab ===
            'cash_vs_mortgage'
                ? 'bg-purple-400/70 text-white'
                : ''}"
            on:click={() => (selectedTab = "cash_vs_mortgage")}
            type="button"
            transition:slide={{ duration: 500 }}
        >
            🔮
            <span
                ><h1 class="text-xs font-bold leading-tight">
                    {$_("navBar.cashCompare")}
                </h1></span
            >
        </button>

        <button
            class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl text-gray-300 transition-all duration-300 ease-in-out hover:bg-gray-800 {selectedTab ===
            'prices'
                ? 'bg-purple-400/70 text-white'
                : ''}"
            on:click={() => (selectedTab = "prices")}
            type="button"
            transition:slide={{ duration: 500 }}
        >
            🏡
            <span
                ><h1 class="text-xs font-bold leading-tight">
                    {$_("navBar.prices")}
                </h1></span
            >
        </button>
    {:else}
        <button
            class="flex flex-col items-center cursor-pointer p-1 flex-1 rounded-xl bg-gradient-to-r from-yellow-400 to-orange-500 animate-pulse text-black transition-all duration-300 ease-in-out hover:bg-gray-800"
            on:click={() => {
                getProFunction();
            }}
            type="button"
            transition:slide={{ duration: 500 }}
        >
            ⭐
            <span><h1 class="text-xs font-bold leading-tight">PRO</h1></span>
        </button>
    {/if}
</div>

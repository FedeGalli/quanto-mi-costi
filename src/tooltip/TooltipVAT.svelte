<script lang="ts">
    import { _ } from "svelte-i18n";
    export let showTooltip: boolean[] = [];
    export let index: number = 0;
</script>

<!-- Tooltip (Overlay) -->
<div
    class="absolute top-0 left-0 w-full h-full bg-[#1e1f25] rounded-2xl shadow-xl p-8 overflow-auto z-20 transition-all transform duration-600 ease-out-[cubic-bezier(0.22, 1, 0.36, 1)]"
    class:opacity-100={showTooltip[index]}
    class:opacity-0={!showTooltip[index]}
    class:-translate-y-500={!showTooltip[index]}
    class:translate-y-0={showTooltip[index]}
>
    <!-- Close Button -->
    <button
        on:click={() => (showTooltip[index] = false)}
        class="absolute top-6 right-6 text-white text-5xl font-bold hover:text-purple-400 transition"
        aria-label="Chiudi tooltip"
    >
        ×
    </button>
    <div class="space-y-6">
        <h2 class="text-4xl font-bold leading-tight">
            {$_("tooltip.vat.title.pay")}
            <span class="text-purple-400">{$_("tooltip.vat.title.vat")}</span
            >{$_("tooltip.vat.title.if")}
        </h2>
        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1. </span>
            {$_("tooltip.vat.paragraph1.title")}
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                {$_("tooltip.vat.paragraph1.bullet1")}
            </li>
            <li>
                {$_("tooltip.vat.paragraph1.bullet2")}
            </li>
            <li>
                {$_("tooltip.vat.paragraph1.bullet3")}
            </li>
        </ul>
    </div>
</div>

<script lang="ts">
    export let showTooltip: boolean[] = [];
    export let index: number = 0;
    import { _ } from "svelte-i18n";
    let tooltipElement: any;

    $: if (showTooltip[index] && tooltipElement) {
        setTimeout(() => {
            // Get the parent container position
            const parent = tooltipElement.parentElement;
            if (parent) {
                const rect = parent.getBoundingClientRect();
                const scrollTop =
                    window.pageYOffset || document.documentElement.scrollTop;

                // Scroll to the parent container
                window.scrollTo({
                    top: scrollTop + rect.top - 100, // 20px padding from top
                    behavior: "smooth",
                });
            }
        }, 100);
    }
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
            {$_("tooltip.firstHouse.title.do")}<span class="text-purple-400"
                >{$_("tooltip.firstHouse.title.first")}</span
            >
            {$_("tooltip.firstHouse.title.if")}
        </h2>
        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1. </span>
            {$_("tooltip.firstHouse.paragraph1.title")}
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                {$_("tooltip.firstHouse.paragraph1.bullet1")}
            </li>
            <li>
                {$_("tooltip.firstHouse.paragraph1.bullet2")}
            </li>
            <li>
                {$_("tooltip.firstHouse.paragraph1.bullet3")}
            </li>
        </ul>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">2. </span>
            {$_("tooltip.firstHouse.paragraph2.title")}
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                {$_("tooltip.firstHouse.paragraph2.bullet1")}
            </li>
            <li>
                {$_("tooltip.firstHouse.paragraph2.bullet2")}
            </li>
        </ul>
    </div>
</div>

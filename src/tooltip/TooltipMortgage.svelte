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
    <div class="space-y-4">
        <h2 class="text-4xl font-bold leading-tight">
            {$_("tooltip.mortgage.title.info")}<span class="text-purple-400"
                >{$_("tooltip.mortgage.title.mortgage")}</span
            >{$_("tooltip.mortgage.title.dots")}
        </h2>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1. </span>
            {$_("tooltip.mortgage.paragraph1.title")}
            <span class="text-base font-medium leading-tight"
                >{$_("tooltip.mortgage.paragraph1.bullet1")}</span
            >
        </p>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">2. </span>
            {$_("tooltip.mortgage.paragraph2.title")}
            <span class="text-base font-medium leading-tight">
                {$_("tooltip.mortgage.paragraph2.bullet1")}</span
            >
        </p>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">3. </span>
            {$_("tooltip.mortgage.paragraph3.title")}
            <span class="text-base font-medium leading-tight">
                {$_("tooltip.mortgage.paragraph3.bullet1")}</span
            >
        </p>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">4. </span>
            {$_("tooltip.mortgage.paragraph4.title")}
        </p>

        <ul class="list-disc pl-6 text-base leading-tight space-y-3">
            <li>
                <span class="font-bold leading-tight"
                    >{$_(
                        "tooltip.mortgage.paragraph4.subParagraph1.title",
                    )}</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph1.bullet1.title",
                            )}</span
                        >{$_(
                            "tooltip.mortgage.paragraph4.subParagraph1.bullet1.content",
                        )}
                    </li>
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph1.bullet2.title",
                            )}</span
                        >{$_(
                            "tooltip.mortgage.paragraph4.subParagraph1.bullet2.content",
                        )}
                    </li>
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph1.bullet3.title",
                            )}</span
                        >{$_(
                            "tooltip.mortgage.paragraph4.subParagraph1.bullet3.content",
                        )}
                    </li>
                </ul>
            </li>
            <li>
                <span class="font-bold leading-tight"
                    >{$_(
                        "tooltip.mortgage.paragraph4.subParagraph2.title",
                    )}</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        {$_(
                            "tooltip.mortgage.paragraph4.subParagraph2.bullet1",
                        )}
                    </li>
                </ul>
            </li>
            <li>
                <span class="font-bold leading-tight"
                    >{$_(
                        "tooltip.mortgage.paragraph4.subParagraph3.title",
                    )}</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph3.bullet1.title",
                            )}</span
                        >{$_(
                            "tooltip.mortgage.paragraph4.subParagraph3.bullet1.content",
                        )}
                    </li>
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph3.bullet2.title",
                            )}</span
                        >{$_(
                            "tooltip.mortgage.paragraph4.subParagraph3.bullet2.content",
                        )}
                    </li>
                </ul>
            </li>
            <li>
                <span class="font-bold leading-tight"
                    >{$_(
                        "tooltip.mortgage.paragraph4.subParagraph4.title",
                    )}</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph4.bullet1.title",
                            )}</span
                        >{$_(
                            "tooltip.mortgage.paragraph4.subParagraph4.bullet1.content",
                        )}
                    </li>
                    <li>
                        <span class="font-bold"
                            >{$_(
                                "tooltip.mortgage.paragraph4.subParagraph4.bullet2.title",
                            )}</span
                        >
                        {$_(
                            "tooltip.mortgage.paragraph4.subParagraph4.bullet2.content",
                        )}
                    </li>
                </ul>
            </li>
        </ul>
    </div>
</div>

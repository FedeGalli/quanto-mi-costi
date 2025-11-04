<script>
    import { _ } from "svelte-i18n";
    import { onMount } from "svelte";
    let showBanner = false;
    let showSettings = false;

    // Cookie preferences
    let preferences = {
        necessary: true, // Always true, can't be disabled
        analytics: true,
        marketing: true,
    };

    onMount(() => {
        // Check if user has already accepted/rejected cookies
        const consent = localStorage.getItem("cookieConsent");
        if (!consent) {
            showBanner = true;
        }
    });

    function acceptAll() {
        preferences = {
            necessary: true,
            analytics: true,
            marketing: true,
        };
        saveCookiePreferences();
    }

    function savePreferences() {
        saveCookiePreferences();
    }

    function saveCookiePreferences() {
        localStorage.setItem("cookieConsent", "true");
        localStorage.setItem("cookiePreferences", JSON.stringify(preferences));
        showBanner = false;
        showSettings = false;

        // Initialize services based on preferences
        if (preferences.analytics) {
            initializeAnalytics();
        }
        if (preferences.marketing) {
            initializeMarketing();
        }
    }

    function toggleSettings() {
        showSettings = !showSettings;
    }

    function initializeAnalytics() {
        //console.log("Analytics initialized");
    }

    function initializeMarketing() {
        //console.log("Marketing cookies initialized");
    }
</script>

{#if showBanner}
    <div
        class="fixed bottom-4 left-4 right-4 md:left-auto md:right-6 md:max-w-md bg-[#1e1f25] text-white rounded-xl shadow-lg border border-[#1e1f25] z-50 animate-fade-in"
    >
        <div class="p-5">
            {#if !showSettings}
                <div class="space-y-4">
                    <p class="text-sm leading-relaxed text-white">
                        {$_("cookie.title")}
                        <a
                            href="/privacy-policy"
                            target="_blank"
                            class="text-purple-400 hover:text-purple-500 underline transition-colors"
                        >
                            {$_("cookie.findMore")}
                        </a>
                    </p>
                    <div class="flex gap-2">
                        <button
                            on:click={toggleSettings}
                            class="flex-1 px-4 py-2 text-sm text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors font-medium"
                        >
                            {$_("cookie.edit")}
                        </button>
                        <button
                            on:click={acceptAll}
                            class="flex-1 px-4 py-2 text-sm bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors font-medium"
                        >
                            {$_("cookie.accept")}
                        </button>
                    </div>
                </div>
            {:else}
                <div class="space-y-4 animate-fade-in">
                    <div class="flex justify-between items-center">
                        <h3 class="text-base font-semibold text-white">
                            {$_("cookie.preferences")}
                        </h3>
                        <button
                            on:click={toggleSettings}
                            class="text-gray-400 hover:text-white transition-colors text-xl leading-none"
                        >
                            ✕
                        </button>
                    </div>

                    <div class="space-y-2">
                        <label
                            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg cursor-not-allowed"
                        >
                            <div>
                                <p class="text-sm font-medium text-gray-700">
                                    {$_("cookie.necessary")}
                                </p>
                                <p class="text-xs text-gray-500">
                                    {$_("cookie.necessaryDesc")}
                                </p>
                            </div>
                            <input
                                type="checkbox"
                                checked
                                disabled
                                class="w-4 h-4 accent-purple-500 opacity-50 cursor-not-allowed"
                            />
                        </label>

                        <label
                            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100 transition-colors"
                        >
                            <div>
                                <p class="text-sm font-medium text-gray-700">
                                    {$_("cookie.analytics")}
                                </p>
                                <p class="text-xs text-gray-500">
                                    {$_("cookie.analyticsDesc")}
                                </p>
                            </div>
                            <input
                                type="checkbox"
                                bind:checked={preferences.analytics}
                                class="w-4 h-4 accent-gray-400 cursor-pointer"
                            />
                        </label>

                        <label
                            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg cursor-pointer hover:bg-gray-100 transition-colors"
                        >
                            <div>
                                <p class="text-sm font-medium text-gray-700">
                                    {$_("cookie.marketing")}
                                </p>
                                <p class="text-xs text-gray-500">
                                    {$_("cookie.marketingDesc")}
                                </p>
                            </div>
                            <input
                                type="checkbox"
                                bind:checked={preferences.marketing}
                                class="w-4 h-4 accent-gray-400 cursor-pointer"
                            />
                        </label>
                    </div>

                    <button
                        on:click={savePreferences}
                        class="w-full px-4 py-2 text-sm bg-purple-500 text-white rounded-lg hover:-purple-600 transition-colors font-medium"
                    >
                        Salva preferenze
                    </button>
                </div>
            {/if}
        </div>
    </div>
{/if}

<style>
    @keyframes fade-in {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
    .animate-fade-in {
        animation: fade-in 0.4s ease-out;
    }
</style>

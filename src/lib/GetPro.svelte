<script lang="ts">
    import { _ } from "svelte-i18n";
    import { fly } from "svelte/transition";
    import { onMount } from "svelte";
    import { push } from "svelte-spa-router";
    import { auth } from "./auth/credentials";
    import {
        user,
        isAuthenticated,
        isLoading,
        initAuthStore,
        logout,
    } from "./auth/auth-store";
    import { onAuthStateChanged } from "firebase/auth";
    // Stato del componente
    let mounted = false;

    // Gestori di eventi
    const goToLogin = () => {
        push("/signin?redirect=" + encodeURIComponent("/checkout"));
    };

    const purchasePro = () => {
        if ($isAuthenticated) {
            goToCheckout();
        } else {
            goToLogin();
        }
    };

    const goToCheckout = () => {
        push("/checkout");
    };

    const tryFree = () => {
        push("/");
    };

    onMount(() => {
        mounted = true;

        initAuthStore();
    });
</script>

<!-- Contenitore Principale -->
<div
    class="min-h-screen bg-gradient-to-b from-purple-400 to-[#1e1f25] flex items-start justify-center p-2 sm:p-6 pt-16 sm:pt-20"
>
    {#if mounted}
        <!-- Back to Home Button -->
        <button
            on:click={() => push("/")}
            class="fixed top-4 left-4 bg-white/20 hover:bg-white/30 text-white px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg z-50 transition-all duration-300 font-medium text-sm sm:text-base backdrop-blur-sm"
        >
            ←
        </button>
        <div
            class="w-full max-w-[1350px] mx-auto"
            in:fly={{ x: -1000, duration: 500, delay: 100 }}
        >
            <!-- Sezione Hero -->
            <div class="text-center mb-8 sm:mb-12">
                <h1
                    class="text-4xl sm:text-5xl lg:text-6xl font-bold text-white mb-4 leading-tight"
                >
                    {$_("getPro.title.house")}<span
                        class="text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 to-orange-500"
                        >{$_("getPro.title.few")}</span
                    >
                </h1>
                <p
                    class="text-xl sm:text-2xl text-purple-100 mb-8 max-w-3xl mx-auto"
                >
                    {$_("getPro.title.subTitle")}
                </p>
            </div>

            <!-- Griglia di Confronto Funzionalità -->
            <div class="flex flex-col lg:flex-row gap-6 sm:gap-8 items-stretch">
                <!-- Tile Funzionalità Standard -->
                <div
                    class="w-full lg:basis-[48%] bg-[#1e1f25] rounded-3xl shadow-2xl p-6 sm:p-8 lg:p-10 border border-gray-500/20 hover:border-gray-400/40 transition-all duration-500 hover:scale-[1.02] group flex flex-col relative"
                    in:fly={{ y: 50, duration: 600, delay: 200 }}
                >
                    <div class="flex-1">
                        <!-- Header con Icona e Titolo -->
                        <div class="flex items-center mb-6">
                            <div
                                class="w-16 h-16 bg-gradient-to-r from-gray-400 to-gray-600 rounded-2xl flex items-center justify-center shadow-lg shadow-gray-400/25 group-hover:animate-pulse"
                            >
                                <svg
                                    class="w-8 h-8 text-white"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                    />
                                </svg>
                            </div>
                            <h3
                                class="text-2xl sm:text-3xl font-bold text-white ml-4"
                            >
                                {$_("getPro.std")}
                            </h3>
                        </div>

                        <!-- Badge Standard -->
                        <div
                            class="absolute top-6 right-6 bg-gradient-to-r from-gray-400 to-gray-600 text-white px-3 py-1 rounded-full text-xs font-bold"
                        >
                            {$_("getPro.free")}
                        </div>

                        <p class="text-gray-300 text-lg mb-6 leading-relaxed">
                            {$_("getPro.freeDesc")}
                        </p>

                        <!-- Vantaggi delle Funzionalità -->
                        <div class="space-y-3 mb-8">
                            {#each [$_("getPro.freeBenefit.row1"), $_("getPro.freeBenefit.row2"), $_("getPro.freeBenefit.row3")] as benefit, i}
                                <div
                                    class="flex items-center text-gray-400"
                                    in:fly={{
                                        x: -20,
                                        duration: 400,
                                        delay: 400 + i * 100,
                                    }}
                                >
                                    <svg
                                        class="w-5 h-5 mr-3 flex-shrink-0"
                                        fill="currentColor"
                                        viewBox="0 0 20 20"
                                    >
                                        <path
                                            fill-rule="evenodd"
                                            d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                            clip-rule="evenodd"
                                        />
                                    </svg>
                                    <span>{benefit}</span>
                                </div>
                            {/each}
                        </div>
                    </div>

                    <!-- Sezione Prezzo -->
                    <div class="mt-auto">
                        <div class="text-center mb-6">
                            <div class="text-3xl font-bold text-white mb-2">
                                {$_("getPro.free")}
                            </div>
                            <div class="text-gray-300 text-sm">
                                {$_("getPro.alwaysFree")}
                            </div>
                        </div>

                        <button
                            on:click={tryFree}
                            class="w-full bg-gray-600 hover:bg-gray-500 text-white font-semibold py-3 px-6 rounded-2xl text-lg transition-all duration-300 hover:scale-105"
                        >
                            {$_("getPro.useFreeVersion")}
                        </button>
                    </div>
                </div>

                <!-- Tile Funzionalità Pro -->
                <div
                    class="w-full lg:basis-[48%] bg-[#1e1f25] rounded-3xl shadow-2xl p-6 sm:p-8 lg:p-10 border border-purple-500/20 hover:border-purple-400/40 transition-all duration-500 hover:scale-[1.02] group flex flex-col relative"
                    in:fly={{ y: 50, duration: 600, delay: 400 }}
                >
                    <div class="flex-1">
                        <!-- Header con Icona e Titolo -->
                        <div class="flex items-center mb-6">
                            <div
                                class="w-16 h-16 bg-gradient-to-r from-purple-500 to-pink-500 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25 group-hover:animate-pulse"
                            >
                                <svg
                                    class="w-8 h-8 text-white"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                                    />
                                </svg>
                            </div>
                            <h3
                                class="text-2xl sm:text-3xl font-bold text-white ml-4"
                            >
                                {$_("getPro.pro")}
                            </h3>
                        </div>

                        <!-- Badge Pro -->
                        <div
                            class="absolute top-6 right-6 bg-gradient-to-r from-yellow-400 to-orange-500 text-black px-3 py-1 rounded-full text-xs font-bold animate-pulse"
                        >
                            {$_("getPro.proCaps")}
                        </div>

                        <p class="text-purple-100 text-lg mb-6 leading-relaxed">
                            {$_("getPro.proDesc")}
                        </p>

                        <!-- Vantaggi delle Funzionalità -->
                        <div class="space-y-3 mb-8">
                            {#each [$_("getPro.proBenefit.row1"), $_("getPro.proBenefit.row2"), $_("getPro.proBenefit.row3"), $_("getPro.proBenefit.row4"), $_("getPro.proBenefit.row5")] as benefit, i}
                                <div
                                    class="flex items-center text-green-400"
                                    in:fly={{
                                        x: -20,
                                        duration: 400,
                                        delay: 600 + i * 100,
                                    }}
                                >
                                    <svg
                                        class="w-5 h-5 mr-3 flex-shrink-0"
                                        fill="currentColor"
                                        viewBox="0 0 20 20"
                                    >
                                        <path
                                            fill-rule="evenodd"
                                            d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                            clip-rule="evenodd"
                                        />
                                    </svg>
                                    <span>{benefit}</span>
                                </div>
                            {/each}
                        </div>
                    </div>

                    <!-- Sezione Prezzo -->
                    <div class="mt-auto">
                        <div class="text-center mb-6">
                            <div class="text-4xl font-bold text-white mb-2">
                                €3,99
                            </div>
                            <div class="text-purple-100 text-sm mb-1">
                                {$_("getPro.vat")}
                            </div>
                            <div class="text-xs text-purple-300">
                                {$_("getPro.access")}
                            </div>
                        </div>

                        <button
                            on:click={purchasePro}
                            class="w-full bg-gradient-to-r from-yellow-400 to-orange-500 hover:from-yellow-500 hover:to-orange-600 text-black font-bold py-4 px-6 rounded-2xl text-lg transition-all duration-300 hover:scale-105 shadow-lg hover:shadow-yellow-400/25 animate-pulse"
                        >
                            {$_("getPro.getProNow")}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
</style>

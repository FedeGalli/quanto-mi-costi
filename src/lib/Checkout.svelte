<script lang="ts">
    import { fly } from "svelte/transition";
    import { onMount } from "svelte";
    import { push } from "svelte-spa-router";
    import { auth, db } from "./auth/credentials";
    import {
        user,
        isAuthenticated,
        isLoading,
        initAuthStore,
        logout,
    } from "./auth/auth-store";
    import { doc, setDoc, getDoc, updateDoc } from "firebase/firestore";

    // Form state
    let mounted = false;
    let isProcessing = false;

    // Form data
    let formData = {
        email: "",
        fullName: "",
        cardNumber: "",
        expiryDate: "",
        cvv: "",
        billingAddress: {
            street: "",
            city: "",
            postalCode: "",
            country: "IT",
        },
    };

    // Validation state
    let errors: Record<string, string> = {};

    // Selected plan (default to Pro)
    let selectedPlan = {
        name: "Pro",
        price: 3.27,
        currency: "€",
        features: [
            "Dashboard eveluzione patrimonio con diverse durate mutuo",
            "Dashboard eveluzione patrimonio con un mutuo o pagando cash",
            "Quanto costano le case nella mia zona",
        ],
    };

    // Event handlers
    const goBack = () => {
        push("/getpro");
    };

    const validateForm = (): boolean => {
        errors = {};

        if (!formData.email) errors.email = "Email richiesta";
        if (!formData.fullName) errors.fullName = "Nome completo richiesto";

        if (!formData.cardNumber) errors.cardNumber = "Numero carta richiesto";
        if (!formData.expiryDate) errors.expiryDate = "Data scadenza richiesta";
        if (!formData.cvv) errors.cvv = "CVV richiesto";

        if (!formData.billingAddress.street)
            errors.street = "Indirizzo richiesto";
        if (!formData.billingAddress.city) errors.city = "Città richiesta";
        if (!formData.billingAddress.postalCode)
            errors.postalCode = "CAP richiesto";

        return Object.keys(errors).length === 0;
    };

    const processPayment = async () => {
        if (!validateForm()) return;

        isProcessing = true;

        try {
            // Simulate payment processing
            await new Promise((resolve) => setTimeout(resolve, 2000));

            // Handle successful payment
            console.log("Payment successful:", formData);
            updateUserPro($user.uid, true);
            push("/"); // Redirect to dashboard or success page
        } catch (error) {
            console.error("Payment failed:", error);
            // Handle payment error
        } finally {
            isProcessing = false;
        }
    };

    async function updateUserPro(uid: string, isPro: boolean) {
        try {
            const userDocRef = doc(db, "users", uid);
            await updateDoc(userDocRef, {
                is_pro: isPro,
            });
            console.log("User pro status updated successfully");
        } catch (error) {
            console.error("Error updating user pro status:", error);
            throw error;
        }
    }

    const formatCardNumber = (value: string) => {
        return value
            .replace(/\s+/g, "")
            .replace(/(.{4})/g, "$1 ")
            .trim();
    };

    const formatExpiryDate = (value: string) => {
        return value.replace(/\D/g, "").replace(/(.{2})(.{2})/, "$1/$2");
    };

    onMount(() => {
        mounted = true;
        initAuthStore();

        if (!$isAuthenticated) {
            push("/signin?redirect=" + encodeURIComponent("/checkout"));
        }
    });
</script>

<!-- Contenitore Principale -->
<div
    class="min-h-screen bg-gradient-to-b from-purple-400 to-[#1e1f25] flex items-start justify-center p-2 sm:p-6 pt-16 sm:pt-20"
>
    {#if mounted}
        <!-- Back to Home Button -->
        <button
            on:click={() => push("/getpro")}
            class="fixed top-4 left-4 bg-white/20 hover:bg-white/30 text-white px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg z-50 transition-all duration-300 font-medium text-sm sm:text-base backdrop-blur-sm"
        >
            ←
        </button>
        <div
            class="w-full max-w-[1350px] mx-auto"
            in:fly={{ x: -1000, duration: 500, delay: 100 }}
        >
            <!-- Payment Layout -->
            <div class="flex flex-col lg:flex-row gap-6 sm:gap-8 items-start">
                <!-- Order Summary -->
                <div
                    class="w-full lg:basis-[35%] bg-[#1e1f25] rounded-3xl shadow-2xl p-6 sm:p-8 border border-purple-500/20 flex flex-col"
                    in:fly={{ y: 50, duration: 600, delay: 200 }}
                >
                    <div class="flex items-center mb-6">
                        <div
                            class="w-16 h-16 bg-gradient-to-r from-purple-500 to-pink-500 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25"
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
                            {selectedPlan.name}
                        </h3>
                    </div>

                    <div class="space-y-3 mb-6">
                        {#each selectedPlan.features as feature, i}
                            <div
                                class="flex items-center text-green-400"
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
                                <span class="text-sm">{feature}</span>
                            </div>
                        {/each}
                    </div>

                    <div class="border-t border-gray-600 pt-6 mt-auto">
                        <div class="flex justify-between items-center mb-2">
                            <span class="text-gray-300"
                                >Piano {selectedPlan.name}</span
                            >
                            <span class="text-white font-semibold"
                                >{selectedPlan.currency}{selectedPlan.price}</span
                            >
                        </div>
                        <div class="flex justify-between items-center mb-4">
                            <span class="text-gray-300">IVA (22%)</span>
                            <span class="text-white font-semibold"
                                >{selectedPlan.currency}{(
                                    selectedPlan.price * 0.22
                                ).toFixed(2)}</span
                            >
                        </div>
                        <div class="border-t border-gray-600 pt-4">
                            <div class="flex justify-between items-center">
                                <span class="text-xl font-bold text-white"
                                    >Totale</span
                                >
                                <span class="text-2xl font-bold text-yellow-400"
                                    >{selectedPlan.currency}{(
                                        selectedPlan.price * 1.22
                                    ).toFixed(2)}</span
                                >
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Payment Form -->
                <div
                    class="w-full lg:basis-[65%] bg-[#1e1f25] rounded-3xl shadow-2xl p-6 sm:p-8 border border-gray-500/20"
                    in:fly={{ y: 50, duration: 600, delay: 400 }}
                >
                    <h2 class="text-2xl font-bold text-white mb-6">
                        Informazioni di pagamento
                    </h2>

                    <form
                        on:submit|preventDefault={processPayment}
                        class="space-y-6"
                    >
                        <!-- Contact Information -->
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                                <label
                                    class="block text-sm font-semibold text-gray-300 mb-2"
                                    >Email</label
                                >
                                <input
                                    type="email"
                                    bind:value={formData.email}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    placeholder="tua@email.com"
                                    required
                                />
                                {#if errors.email}
                                    <p class="text-red-400 text-sm mt-1">
                                        {errors.email}
                                    </p>
                                {/if}
                            </div>
                            <div>
                                <label
                                    class="block text-sm font-semibold text-gray-300 mb-2"
                                    >Nome Completo</label
                                >
                                <input
                                    type="text"
                                    bind:value={formData.fullName}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    placeholder="Mario Rossi"
                                    required
                                />
                                {#if errors.fullName}
                                    <p class="text-red-400 text-sm mt-1">
                                        {errors.fullName}
                                    </p>
                                {/if}
                            </div>
                        </div>
                        <!-- Card Information -->
                        <div class="space-y-4">
                            <div>
                                <label
                                    class="block text-sm font-semibold text-gray-300 mb-2"
                                    >Numero Carta</label
                                >
                                <input
                                    type="text"
                                    bind:value={formData.cardNumber}
                                    on:input={(e) =>
                                        (formData.cardNumber = formatCardNumber(
                                            e.target.value,
                                        ))}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    placeholder="1234 5678 9012 3456"
                                    maxlength="19"
                                    required
                                />
                                {#if errors.cardNumber}
                                    <p class="text-red-400 text-sm mt-1">
                                        {errors.cardNumber}
                                    </p>
                                {/if}
                            </div>

                            <div class="grid grid-cols-2 gap-4">
                                <div>
                                    <label
                                        class="block text-sm font-semibold text-gray-300 mb-2"
                                        >Data Scadenza</label
                                    >
                                    <input
                                        type="text"
                                        bind:value={formData.expiryDate}
                                        on:input={(e) =>
                                            (formData.expiryDate =
                                                formatExpiryDate(
                                                    e.target.value,
                                                ))}
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        placeholder="MM/AA"
                                        maxlength="5"
                                        required
                                    />
                                    {#if errors.expiryDate}
                                        <p class="text-red-400 text-sm mt-1">
                                            {errors.expiryDate}
                                        </p>
                                    {/if}
                                </div>

                                <div>
                                    <label
                                        class="block text-sm font-semibold text-gray-300 mb-2"
                                        >CVV</label
                                    >
                                    <input
                                        type="text"
                                        bind:value={formData.cvv}
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        placeholder="123"
                                        maxlength="4"
                                        required
                                    />
                                    {#if errors.cvv}
                                        <p class="text-red-400 text-sm mt-1">
                                            {errors.cvv}
                                        </p>
                                    {/if}
                                </div>
                            </div>
                        </div>

                        <!-- Billing Address -->
                        <div class="space-y-4">
                            <h3 class="text-lg font-semibold text-white">
                                Indirizzo di Fatturazione
                            </h3>

                            <div>
                                <label
                                    class="block text-sm font-semibold text-gray-300 mb-2"
                                    >Indirizzo</label
                                >
                                <input
                                    type="text"
                                    bind:value={formData.billingAddress.street}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    placeholder="Via Roma, 123"
                                    required
                                />
                                {#if errors.street}
                                    <p class="text-red-400 text-sm mt-1">
                                        {errors.street}
                                    </p>
                                {/if}
                            </div>

                            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                                <div>
                                    <label
                                        class="block text-sm font-semibold text-gray-300 mb-2"
                                        >Città</label
                                    >
                                    <input
                                        type="text"
                                        bind:value={
                                            formData.billingAddress.city
                                        }
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        placeholder="Milano"
                                        required
                                    />
                                    {#if errors.city}
                                        <p class="text-red-400 text-sm mt-1">
                                            {errors.city}
                                        </p>
                                    {/if}
                                </div>

                                <div>
                                    <label
                                        class="block text-sm font-semibold text-gray-300 mb-2"
                                        >CAP</label
                                    >
                                    <input
                                        type="text"
                                        bind:value={
                                            formData.billingAddress.postalCode
                                        }
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        placeholder="20100"
                                        required
                                    />
                                    {#if errors.postalCode}
                                        <p class="text-red-400 text-sm mt-1">
                                            {errors.postalCode}
                                        </p>
                                    {/if}
                                </div>

                                <div>
                                    <label
                                        class="block text-sm font-semibold text-gray-300 mb-2"
                                        >Paese</label
                                    >
                                    <select
                                        bind:value={
                                            formData.billingAddress.country
                                        }
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        required
                                    >
                                        <option value="IT">Italia</option>
                                        <option value="FR">Francia</option>
                                        <option value="DE">Germania</option>
                                        <option value="ES">Spagna</option>
                                    </select>
                                </div>
                            </div>
                        </div>

                        <!-- Action Buttons -->
                        <div class="flex flex-col sm:flex-row gap-4 pt-6">
                            <button
                                type="submit"
                                disabled={isProcessing}
                                class="flex-1 bg-gradient-to-r from-yellow-400 to-orange-500 hover:from-yellow-500 hover:to-orange-600 disabled:from-gray-500 disabled:to-gray-600 text-black font-bold py-3 px-6 rounded-2xl text-lg transition-all duration-300 hover:scale-105 shadow-lg hover:shadow-yellow-400/25 disabled:cursor-not-allowed disabled:hover:scale-100"
                            >
                                {#if isProcessing}
                                    <div
                                        class="flex items-center justify-center"
                                    >
                                        <svg
                                            class="animate-spin -ml-1 mr-3 h-5 w-5 text-black"
                                            xmlns="http://www.w3.org/2000/svg"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                        >
                                            <circle
                                                class="opacity-25"
                                                cx="12"
                                                cy="12"
                                                r="10"
                                                stroke="currentColor"
                                                stroke-width="4"
                                            ></circle>
                                            <path
                                                class="opacity-75"
                                                fill="currentColor"
                                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                            ></path>
                                        </svg>
                                        Elaborazione...
                                    </div>
                                {:else}
                                    Procedi al pagamento
                                {/if}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    /* Additional custom styles if needed */
</style>

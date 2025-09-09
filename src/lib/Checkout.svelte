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

    // Stripe configuration
    const STRIPE_PUBLISHABLE_KEY =
        "pk_test_51S54jECPlqlBCAOphwqE8Z9bhP6crRGY6xQWXEgL5m9XWqPiOQhbOgLtM74HpneTrGYhxVqHaobEfYPOYV3K6hJB009OULWyrH";

    // Stripe variables
    let stripe: any = null;
    let elements: any = null;
    let cardElement: any = null;
    let stripeLoaded = false;
    let showSuccessPopup = false;

    // Form state
    let mounted = false;
    let isProcessing = false;

    // Form data (simplified since Stripe handles card details)
    let formData = {
        email: "",
        fullName: "",
        billingAddress: {
            street: "",
            city: "",
            postalCode: "",
            country: "IT",
        },
    };

    // Validation state
    let errors: Record<string, string> = {};
    let cardErrors = "";

    // Selected plan (default to Pro)
    let selectedPlan = {
        name: "Pro",
        price: 3.27,
        currency: "€",
        features: [
            "Conviene fare un muto o pagare cash?",
            "Che durata del mutuo mi conviene?",
            "Quanto costano le case nella mia zona?",
            "Quali sono i volumi di compravendita nella mia zona?",
            "Nel passato come erano i prezzi e volumi di compravendita?",
        ],
    };

    // Initialize Stripe
    const initStripe = async () => {
        try {
            // Load Stripe.js
            const script = document.createElement("script");
            script.src = "https://js.stripe.com/v3/";
            script.onload = () => {
                stripe = (window as any).Stripe(STRIPE_PUBLISHABLE_KEY);

                // Create Elements instance
                elements = stripe.elements({
                    appearance: {
                        theme: "night",
                        variables: {
                            colorPrimary: "#8b5cf6",
                            colorBackground: "#1f2937",
                            colorText: "#ffffff",
                            colorDanger: "#ef4444",
                            fontFamily: "system-ui, sans-serif",
                            spacingUnit: "4px",
                            borderRadius: "12px",
                        },
                    },
                });

                // Create card element
                cardElement = elements.create("card", {
                    style: {
                        base: {
                            fontSize: "16px",
                            color: "#ffffff",
                            "::placeholder": {
                                color: "#9ca3af",
                            },
                        },
                        invalid: {
                            color: "#ef4444",
                        },
                    },
                    disableLink: true,
                });

                // Mount card element
                cardElement.mount("#card-element");

                // Listen for real-time validation errors
                cardElement.on("change", (event: any) => {
                    if (event.error) {
                        cardErrors = event.error.message;
                    } else {
                        cardErrors = "";
                    }
                });

                stripeLoaded = true;
            };
            document.head.appendChild(script);
        } catch (error) {
            console.error("Error loading Stripe:", error);
        }
    };

    const processPayment = async () => {
        if (!validateForm() || !stripe || !cardElement) {
            return;
        }

        isProcessing = true;
        cardErrors = "";

        try {
            // Step 1: Create Payment Intent on backend
            const totalAmount = selectedPlan.price * 1.22; // Include VAT
            const amountInCents = Math.round(totalAmount * 100);

            console.log("Creating payment intent with amount:", amountInCents);

            const paymentIntentResponse = await fetch(
                "http://localhost:8080/api/create-payment-intent",
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        amount: amountInCents,
                        currency: "eur",
                        email: formData.email,
                        name: formData.fullName,
                        address: {
                            line1: formData.billingAddress.street,
                            city: formData.billingAddress.city,
                            postal_code: formData.billingAddress.postalCode,
                            country: formData.billingAddress.country,
                        },
                        uid: $user.uid,
                    }),
                },
            );

            if (!paymentIntentResponse.ok) {
                const errorData = await paymentIntentResponse.json();
                throw new Error(
                    errorData.message || "Failed to create payment intent",
                );
            }

            const paymentIntentData = await paymentIntentResponse.json();
            console.log("Payment intent created:", paymentIntentData);

            // Step 2: Confirm payment with the card element
            const { error, paymentIntent } = await stripe.confirmCardPayment(
                paymentIntentData.client_secret,
                {
                    payment_method: {
                        card: cardElement,
                        billing_details: {
                            name: formData.fullName,
                            email: formData.email,
                            address: {
                                line1: formData.billingAddress.street,
                                city: formData.billingAddress.city,
                                postal_code: formData.billingAddress.postalCode,
                                country: formData.billingAddress.country,
                            },
                        },
                    },
                },
            );

            if (error) {
                console.error("Payment failed:", error);
                cardErrors =
                    error.message ||
                    "Si è verificato un errore durante il pagamento. Riprova.";
                return;
            }
            // Step 3: Notify backend of successful payment
            if (paymentIntent.status === "succeeded") {
                const confirmResponse = await fetch(
                    "http://localhost:8080/api/confirm-payment",
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                            payment_intent_id: paymentIntent.id,
                            uid: $user.uid,
                        }),
                    },
                );
                if (!confirmResponse.ok) {
                    const errorData = await confirmResponse.json();
                    throw new Error(
                        errorData.message ||
                            "Failed to confirm payment with server",
                    );
                }

                const confirmData = await confirmResponse.json();
                console.log("Payment confirmed with server:", confirmData);

                if (confirmData.success) {
                    console.log("PRO ACCESS ACTIVATED!");
                    await showProActivatedPopup();
                } else {
                    throw new Error(
                        confirmData.message || "Failed to activate pro access",
                    );
                }
            } else {
                throw new Error(
                    `Payment not completed. Status: ${paymentIntent.status}`,
                );
            }
        } catch (error) {
            console.error("Payment failed:", error.message);
            cardErrors =
                error.message ||
                "Si è verificato un errore durante il pagamento. Riprova.";
        } finally {
            isProcessing = false;
        }
    };

    onMount(() => {
        mounted = true;
        initAuthStore();
        initStripe();

        if (!$isAuthenticated) {
            push("/signin?redirect=" + encodeURIComponent("/checkout"));
        }
        if ($user.pro) {
            push("/");
        }
    });

    const showProActivatedPopup = async () => {
        showSuccessPopup = true;
        console.log(showSuccessPopup);

        // Auto-close after 4 seconds and redirect
        setTimeout(() => {
            showSuccessPopup = false;
            push("/");
        }, 4000);
    };

    // Email validation
    const isValidEmail = (email: string): boolean => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    };

    // Italian postal code validation (5 digits)
    const isValidPostalCode = (postalCode: string): boolean => {
        const postalCodeRegex = /^[0-9]{5}$/;
        return postalCodeRegex.test(postalCode);
    };

    // Name validation (only letters, spaces, apostrophes, hyphens, accents)
    const isValidName = (name: string): boolean => {
        const nameRegex = /^[a-zA-ZÀ-ÿ\s'-]{2,50}$/;
        return nameRegex.test(name.trim());
    };

    // Address validation (letters, numbers, spaces, common punctuation)
    const isValidAddress = (address: string): boolean => {
        const addressRegex = /^[a-zA-Z0-9À-ÿ\s,.-]{3,100}$/;
        return addressRegex.test(address.trim());
    };

    // City validation (only letters, spaces, apostrophes, hyphens, accents)
    const isValidCity = (city: string): boolean => {
        const cityRegex = /^[a-zA-ZÀ-ÿ\s'-]{2,50}$/;
        return cityRegex.test(city.trim());
    };

    // Input handlers for real-time validation and formatting
    const handleEmailInput = (event: Event) => {
        const target = event.target as HTMLInputElement;
        formData.email = target.value.trim().toLowerCase();

        if (formData.email && !isValidEmail(formData.email)) {
            errors.email = "Inserisci un indirizzo email valido";
        } else {
            errors.email = "";
        }
    };

    const handleNameInput = (event: Event) => {
        const target = event.target as HTMLInputElement;
        let value = target.value;

        // Remove any numbers or special characters (except spaces, apostrophes, hyphens)
        value = value.replace(/[^a-zA-ZÀ-ÿ\s'-]/g, "");

        // Capitalize first letter of each word
        value = value.replace(/\b\w/g, (l) => l.toUpperCase());

        formData.fullName = value;
        target.value = value;

        if (formData.fullName && !isValidName(formData.fullName)) {
            errors.fullName =
                "Il nome deve contenere solo lettere (2-50 caratteri)";
        } else if (
            formData.fullName &&
            formData.fullName.trim().split(" ").length < 2
        ) {
            errors.fullName = "Inserisci nome e cognome";
        } else {
            errors.fullName = "";
        }
    };

    const handleAddressInput = (event: Event) => {
        const target = event.target as HTMLInputElement;
        formData.billingAddress.street = target.value;

        if (
            formData.billingAddress.street &&
            !isValidAddress(formData.billingAddress.street)
        ) {
            errors.street =
                "Indirizzo non valido (3-100 caratteri, solo lettere, numeri e punteggiatura comune)";
        } else {
            errors.street = "";
        }
    };

    const handleCityInput = (event: Event) => {
        const target = event.target as HTMLInputElement;
        let value = target.value;

        // Remove numbers and special characters (except spaces, apostrophes, hyphens)
        value = value.replace(/[^a-zA-ZÀ-ÿ\s'-]/g, "");

        // Capitalize first letter of each word
        value = value.replace(/\b\w/g, (l) => l.toUpperCase());

        formData.billingAddress.city = value;
        target.value = value;

        if (
            formData.billingAddress.city &&
            !isValidCity(formData.billingAddress.city)
        ) {
            errors.city =
                "La città deve contenere solo lettere (2-50 caratteri)";
        } else {
            errors.city = "";
        }
    };

    const handlePostalCodeInput = (event: Event) => {
        const target = event.target as HTMLInputElement;
        // Remove any non-digit characters
        let value = target.value.replace(/[^0-9]/g, "");

        // Limit to 5 digits
        if (value.length > 5) {
            value = value.slice(0, 5);
        }

        formData.billingAddress.postalCode = value;
        target.value = value;

        if (
            formData.billingAddress.postalCode &&
            !isValidPostalCode(formData.billingAddress.postalCode)
        ) {
            errors.postalCode = "Il CAP deve essere di esattamente 5 cifre";
        } else {
            errors.postalCode = "";
        }
    };

    const handlePostalCodeKeypress = (event: KeyboardEvent) => {
        // Prevent non-digit characters from being entered
        const char = String.fromCharCode(event.which || event.keyCode);
        if (
            !/[0-9]/.test(char) &&
            !["Backspace", "Delete", "ArrowLeft", "ArrowRight", "Tab"].includes(
                event.key,
            )
        ) {
            event.preventDefault();
        }
    };

    // Enhanced form validation
    const validateForm = (): boolean => {
        errors = {};
        let isValid = true;

        // Email validation
        if (!formData.email) {
            errors.email = "Email richiesta";
            isValid = false;
        } else if (!isValidEmail(formData.email)) {
            errors.email = "Inserisci un indirizzo email valido";
            isValid = false;
        }

        // Full name validation
        if (!formData.fullName) {
            errors.fullName = "Nome completo richiesto";
            isValid = false;
        } else if (!isValidName(formData.fullName)) {
            errors.fullName =
                "Il nome deve contenere solo lettere (2-50 caratteri)";
            isValid = false;
        } else if (formData.fullName.trim().split(" ").length < 2) {
            errors.fullName = "Inserisci nome e cognome";
            isValid = false;
        }

        // Street address validation
        if (!formData.billingAddress.street) {
            errors.street = "Indirizzo richiesto";
            isValid = false;
        } else if (!isValidAddress(formData.billingAddress.street)) {
            errors.street = "Indirizzo non valido";
            isValid = false;
        }

        // City validation
        if (!formData.billingAddress.city) {
            errors.city = "Città richiesta";
            isValid = false;
        } else if (!isValidCity(formData.billingAddress.city)) {
            errors.city = "Nome città non valido";
            isValid = false;
        }

        // Postal code validation
        if (!formData.billingAddress.postalCode) {
            errors.postalCode = "CAP richiesto";
            isValid = false;
        } else if (!isValidPostalCode(formData.billingAddress.postalCode)) {
            errors.postalCode = "Il CAP deve essere di esattamente 5 cifre";
            isValid = false;
        }

        return isValid;
    };

    // Validation on blur events (when user confirms input)
    const handleFieldBlur = (field: string) => {
        switch (field) {
            case "email":
                if (!formData.email) {
                    errors.email = "Email richiesta";
                } else if (!isValidEmail(formData.email)) {
                    errors.email = "Inserisci un indirizzo email valido";
                } else {
                    errors.email = "";
                }
                break;
            case "fullName":
                if (!formData.fullName) {
                    errors.fullName = "Nome completo richiesto";
                } else if (!isValidName(formData.fullName)) {
                    errors.fullName =
                        "Il nome deve contenere solo lettere (2-50 caratteri)";
                } else if (formData.fullName.trim().split(" ").length < 2) {
                    errors.fullName = "Inserisci nome e cognome";
                } else {
                    errors.fullName = "";
                }
                break;
            case "street":
                if (!formData.billingAddress.street) {
                    errors.street = "Indirizzo richiesto";
                } else if (!isValidAddress(formData.billingAddress.street)) {
                    errors.street = "Indirizzo non valido (3-100 caratteri)";
                } else {
                    errors.street = "";
                }
                break;
            case "city":
                if (!formData.billingAddress.city) {
                    errors.city = "Città richiesta";
                } else if (!isValidCity(formData.billingAddress.city)) {
                    errors.city = "Nome città non valido (solo lettere)";
                } else {
                    errors.city = "";
                }
                break;
            case "postalCode":
                if (!formData.billingAddress.postalCode) {
                    errors.postalCode = "CAP richiesto";
                } else if (
                    !isValidPostalCode(formData.billingAddress.postalCode)
                ) {
                    errors.postalCode =
                        "Il CAP deve essere di esattamente 5 cifre";
                } else {
                    errors.postalCode = "";
                }
                break;
        }
        errors = { ...errors }; // Trigger reactivity
    };
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
                                    on:input={handleEmailInput}
                                    on:blur={() => handleFieldBlur("email")}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    class:border-red-500={errors.email}
                                    class:focus:border-red-500={errors.email}
                                    placeholder="tua@email.com"
                                    autocomplete="email"
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
                                    on:input={handleNameInput}
                                    on:blur={() => handleFieldBlur("fullName")}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    class:border-red-500={errors.fullName}
                                    class:focus:border-red-500={errors.fullName}
                                    placeholder="Mario Rossi"
                                    autocomplete="name"
                                    maxlength="50"
                                    required
                                />
                                {#if errors.fullName}
                                    <p class="text-red-400 text-sm mt-1">
                                        {errors.fullName}
                                    </p>
                                {/if}
                            </div>
                        </div>

                        <!-- Stripe Card Element -->
                        <div class="space-y-4">
                            <label
                                class="block text-sm font-semibold text-gray-300 mb-2"
                            >
                                Informazioni Carta di Credito
                            </label>
                            <div
                                id="card-element"
                                class="w-full p-4 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus-within:border-purple-500 focus-within:ring-2 focus-within:ring-purple-500/20 transition-all duration-300 min-h-[50px]"
                            >
                                {#if !stripeLoaded}
                                    <div
                                        class="flex items-center justify-center py-2"
                                    >
                                        <div
                                            class="animate-spin rounded-full h-5 w-5 border-b-2 border-purple-500"
                                        ></div>
                                        <span class="ml-2 text-gray-400"
                                            >Caricamento Stripe...</span
                                        >
                                    </div>
                                {/if}
                            </div>
                            {#if cardErrors}
                                <p class="text-red-400 text-sm mt-1">
                                    {cardErrors}
                                </p>
                            {/if}
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
                                    on:input={handleAddressInput}
                                    on:blur={() => handleFieldBlur("street")}
                                    class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                    class:border-red-500={errors.street}
                                    class:focus:border-red-500={errors.street}
                                    placeholder="Via Roma, 123"
                                    autocomplete="street-address"
                                    maxlength="100"
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
                                        on:input={handleCityInput}
                                        on:blur={() => handleFieldBlur("city")}
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        class:border-red-500={errors.city}
                                        class:focus:border-red-500={errors.city}
                                        placeholder="Milano"
                                        autocomplete="address-level2"
                                        maxlength="50"
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
                                        on:input={handlePostalCodeInput}
                                        on:keypress={handlePostalCodeKeypress}
                                        on:blur={() =>
                                            handleFieldBlur("postalCode")}
                                        class="w-full p-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition-all duration-300"
                                        class:border-red-500={errors.postalCode}
                                        class:focus:border-red-500={errors.postalCode}
                                        placeholder="20100"
                                        autocomplete="postal-code"
                                        maxlength="5"
                                        inputmode="numeric"
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
                                        autocomplete="country"
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
                                disabled={isProcessing || !stripeLoaded}
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
                                {:else if !stripeLoaded}
                                    <div
                                        class="flex items-center justify-center"
                                    >
                                        <div
                                            class="animate-spin rounded-full h-5 w-5 border-b-2 border-black mr-2"
                                        ></div>
                                        Caricamento...
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

{#if showSuccessPopup}
    <!-- Backdrop -->
    <div
        class="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4"
        in:fly={{ y: 50, duration: 500 }}
        out:fly={{ y: -50, duration: 300 }}
    >
        <!-- Popup Content -->
        <div
            class="bg-gradient-to-br from-purple-900 to-indigo-900 rounded-3xl p-8 max-w-md w-full text-center border border-purple-500/30 shadow-2xl relative overflow-hidden"
            in:fly={{ duration: 500, delay: 200 }}
        >
            <!-- Animated Background Elements -->
            <div class="absolute inset-0 overflow-hidden">
                <div
                    class="absolute -top-10 -right-10 w-32 h-32 bg-gradient-to-r from-yellow-400/20 to-orange-500/20 rounded-full animate-pulse"
                ></div>
                <div
                    class="absolute -bottom-10 -left-10 w-24 h-24 bg-gradient-to-r from-purple-500/20 to-pink-500/20 rounded-full animate-pulse"
                    style="animation-delay: 1s;"
                ></div>
            </div>

            <!-- Success Icon -->
            <div class="relative z-10 mb-6">
                <div
                    class="w-20 h-20 bg-gradient-to-r from-yellow-400 to-orange-500 rounded-full flex items-center justify-center mx-auto shadow-lg animate-bounce"
                >
                    <svg
                        class="w-10 h-10 text-white"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="3"
                            d="M5 13l4 4L19 7"
                        ></path>
                    </svg>
                </div>
            </div>

            <!-- Success Text -->
            <div class="relative z-10">
                <h2 class="text-3xl font-bold text-white mb-4">
                    🎉 Pro Attivato!
                </h2>
                <p class="text-gray-300 mb-6 text-lg">
                    Il pagamento è stato completato con successo!<br />
                    Ora hai accesso a tutte le funzionalità Pro.
                </p>

                <!-- Features List -->
                <div class="space-y-2 mb-6">
                    <div
                        class="flex items-center justify-center text-green-400"
                    >
                        <svg
                            class="w-5 h-5 mr-2"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                clip-rule="evenodd"
                            ></path>
                        </svg>
                        <span class="text-sm"
                            >Accesso illimitato a tutte le funzionalità</span
                        >
                    </div>
                    <div
                        class="flex items-center justify-center text-green-400"
                    >
                        <svg
                            class="w-5 h-5 mr-2"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                clip-rule="evenodd"
                            ></path>
                        </svg>
                        <span class="text-sm"
                            >Analisi avanzate del mercato immobiliare</span
                        >
                    </div>
                    <div
                        class="flex items-center justify-center text-green-400"
                    >
                        <svg
                            class="w-5 h-5 mr-2"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                clip-rule="evenodd"
                            ></path>
                        </svg>
                        <span class="text-sm">Supporto prioritario</span>
                    </div>
                </div>

                <!-- Loading bar -->
                <div class="w-full bg-gray-700 rounded-full h-2 mb-4">
                    <div
                        class="bg-gradient-to-r from-yellow-400 to-orange-500 h-2 rounded-full animate-pulse"
                        style="width: 100%; animation: loadingBar 4s linear;"
                    ></div>
                </div>

                <p class="text-gray-400 text-sm">
                    Verrai reindirizzato alla dashboard...
                </p>
            </div>
        </div>
    </div>
{/if}

<style>
    @keyframes loadingBar {
        0% {
            width: 0%;
        }
        100% {
            width: 100%;
        }
    }
</style>

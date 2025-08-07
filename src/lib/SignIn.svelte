<!-- Login Component - Updated Script Section with Enhanced Google Login -->
<script lang="ts">
    import { slide, fly } from "svelte/transition";
    import { doc, setDoc, getDoc } from "firebase/firestore";
    import { onMount } from "svelte";
    import { push, location, querystring } from "svelte-spa-router";
    import {
        createUserWithEmailAndPassword,
        signInWithEmailAndPassword,
        updateProfile,
        GoogleAuthProvider,
        FacebookAuthProvider,
        signInWithPopup,
        signInWithRedirect,
        getRedirectResult,
        sendEmailVerification,
        sendPasswordResetEmail,
        onAuthStateChanged,
    } from "firebase/auth";
    import type { AuthError } from "firebase/auth";
    import { auth, db } from "./auth/credentials"; // Adjust path as needed
    import { user, isAuthenticated, initAuthStore } from "./auth/auth-store"; // Import auth store

    // Redirect
    let redirectUrl = "/"; // Default redirect to homepage

    // Auth state
    let isSignUp = false;
    let isLoading = false;
    let showErrorPopup = false;
    let showSuccessPopup = false;
    let errorMessage = "";
    let successMessage = "";
    let showForgotPassword = false;
    let showEmailVerificationMessage = false;

    // Form data
    let email = "";
    let password = "";
    let confirmPassword = "";
    let firstName = "";
    let lastName = "";
    let acceptTerms = false;
    let resetEmail = "";

    // Validation states
    let emailValid = true;
    let passwordValid = true;
    let confirmPasswordValid = true;
    let resetEmailValid = true;

    // Toggle between sign in and sign up
    function toggleMode() {
        isSignUp = !isSignUp;
        showForgotPassword = false;
        showEmailVerificationMessage = false;
        clearForm();
    }

    // Toggle forgot password view
    function toggleForgotPassword() {
        showForgotPassword = !showForgotPassword;
        showEmailVerificationMessage = false;
        clearForm();
    }

    // Clear form data
    function clearForm() {
        email = "";
        password = "";
        confirmPassword = "";
        firstName = "";
        lastName = "";
        acceptTerms = false;
        resetEmail = "";
        emailValid = true;
        passwordValid = true;
        confirmPasswordValid = true;
        resetEmailValid = true;
    }

    // Validate email
    function validateEmail(email: string): boolean {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return emailRegex.test(email);
    }

    // Validate password
    function validatePassword(password: string): boolean {
        return password.length >= 8;
    }

    // Handle Firebase Auth errors
    function handleFirebaseError(error: AuthError) {
        switch (error.code) {
            case "auth/email-already-in-use":
                return "Questo indirizzo email è già registrato";
            case "auth/invalid-email":
                return "Indirizzo email non valido";
            case "auth/operation-not-allowed":
                return "Operazione non consentita";
            case "auth/weak-password":
                return "La password è troppo debole";
            case "auth/user-disabled":
                return "Account disabilitato";
            case "auth/user-not-found":
                return "Utente non trovato";
            case "auth/wrong-password":
                return "Password errata";
            case "auth/invalid-credential":
                return "Credenziali non valide";
            case "auth/too-many-requests":
                return "Troppi tentativi. Riprova più tardi";
            case "auth/email-not-verified":
                return "Email non verificata. Controlla la tua casella di posta";
            case "auth/popup-closed-by-user":
                return "Popup chiuso dall'utente";
            case "auth/popup-blocked":
                return "Popup bloccato dal browser. Abilita i popup per questo sito";
            case "auth/cancelled-popup-request":
                return "Richiesta popup cancellata";
            case "auth/account-exists-with-different-credential":
                return "Esiste già un account con questo indirizzo email ma con un metodo di accesso diverso";
            default:
                return "Errore di autenticazione. Riprova più tardi";
        }
    }

    // Handle form submission
    async function handleSubmit() {
        // Reset validation states
        emailValid = validateEmail(email);
        passwordValid = validatePassword(password);

        if (isSignUp) {
            confirmPasswordValid = password === confirmPassword;
        }

        // Check if form is valid
        if (
            !emailValid ||
            !passwordValid ||
            (isSignUp && !confirmPasswordValid)
        ) {
            return;
        }

        if (isSignUp && !acceptTerms) {
            showError("Devi accettare i termini e condizioni");
            return;
        }

        isLoading = true;

        try {
            if (isSignUp) {
                // Create new user
                const userCredential = await createUserWithEmailAndPassword(
                    auth,
                    email,
                    password,
                );
                const firebaseUser = userCredential.user;

                // Update user profile with name
                await updateProfile(firebaseUser, {
                    displayName: `${firstName} ${lastName}`,
                });

                // Send email verification
                await sendEmailVerification(firebaseUser);

                showSuccess(
                    "Registrazione completata! Controlla la tua email per verificare l'account.",
                );
                showEmailVerificationMessage = true;
            } else {
                // Sign in existing user
                const userCredential = await signInWithEmailAndPassword(
                    auth,
                    email,
                    password,
                );
                const firebaseUser = userCredential.user;

                // Check if email is verified
                if (!firebaseUser.emailVerified) {
                    showError(
                        "Email non verificata. Controlla la tua casella di posta per il link di verifica.",
                    );
                    showEmailVerificationMessage = true;
                    return;
                }

                showSuccess("Accesso effettuato con successo!");

                // The auth store will handle the user state automatically
                // Redirect to home page after a short delay
                setTimeout(() => {
                    push(redirectUrl);
                }, 1500);
            }
        } catch (error) {
            console.error("Authentication error:", error);
            showError(handleFirebaseError(error as AuthError));
        } finally {
            isLoading = false;
        }
    }

    // Handle password reset
    async function handlePasswordReset() {
        resetEmailValid = validateEmail(resetEmail);

        if (!resetEmailValid) {
            return;
        }

        isLoading = true;

        try {
            await sendPasswordResetEmail(auth, resetEmail);
            showSuccess(
                "Email di reset password inviata! Controlla la tua casella di posta.",
            );
            showForgotPassword = false;
            resetEmail = "";
        } catch (error) {
            console.error("Password reset error:", error);
            showError(handleFirebaseError(error as AuthError));
        } finally {
            isLoading = false;
        }
    }

    // Resend email verification
    async function resendEmailVerification() {
        if (!auth.currentUser) {
            showError("Nessun utente trovato. Riprova la registrazione.");
            return;
        }

        isLoading = true;

        try {
            await sendEmailVerification(auth.currentUser);
            showSuccess("Email di verifica inviata nuovamente!");
        } catch (error) {
            console.error("Resend verification error:", error);
            showError(handleFirebaseError(error as AuthError));
        } finally {
            isLoading = false;
        }
    }

    async function createUserDocument(uid: string) {
        try {
            const userDocRef = doc(db, "users", uid);

            // Check if document already exists
            const userDoc = await getDoc(userDocRef);

            if (!userDoc.exists()) {
                // Create new user document
                await setDoc(userDocRef, {
                    uid: uid,
                    is_pro: false,
                });
                console.log("User document created successfully");
            } else {
                console.log("User document already exists");
            }
        } catch (error) {
            console.error("Error creating user document:", error);
            // Don't throw error to avoid disrupting the login flow
        }
    }

    // Check email verification status
    async function checkEmailVerification() {
        if (!auth.currentUser) {
            return;
        }

        isLoading = true;

        try {
            // Reload user to get latest verification status
            await auth.currentUser.reload();

            if (auth.currentUser.emailVerified) {
                showSuccess("Email verificata! Benvenuto!");
                showEmailVerificationMessage = false;

                createUserDocument(auth.currentUser.uid);
                setTimeout(() => {
                    push(redirectUrl);
                }, 1500);
            } else {
                showError(
                    "Email non ancora verificata. Controlla la tua casella di posta.",
                );
            }
        } catch (error) {
            console.error("Check verification error:", error);
            showError("Errore durante la verifica. Riprova.");
        } finally {
            isLoading = false;
        }
    }

    // Show error popup
    function showError(message: string) {
        errorMessage = message;
        showErrorPopup = true;
        setTimeout(() => {
            showErrorPopup = false;
        }, 5000);
    }

    // Show success popup
    function showSuccess(message: string) {
        successMessage = message;
        showSuccessPopup = true;
        setTimeout(() => {
            showSuccessPopup = false;
        }, 5000);
    }

    // Enhanced Google login with better error handling and configuration
    async function handleGoogleLogin() {
        isLoading = true;
        try {
            // Create Google provider with additional configuration
            const provider = new GoogleAuthProvider();

            // Optional: Add custom parameters
            provider.addScope("email");
            provider.addScope("profile");

            // Optional: Set custom parameters for better UX
            provider.setCustomParameters({
                prompt: "select_account",
            });

            // Try popup first, fall back to redirect on mobile or popup blocked
            let result;
            try {
                result = await signInWithPopup(auth, provider);
            } catch (popupError: any) {
                // If popup fails (blocked, mobile, etc.), use redirect
                if (
                    popupError.code === "auth/popup-blocked" ||
                    popupError.code === "auth/popup-closed-by-user" ||
                    /mobile/i.test(navigator.userAgent)
                ) {
                    console.log(
                        "Popup blocked or mobile detected, using redirect method",
                    );
                    await signInWithRedirect(auth, provider);
                    // The redirect will handle the rest, so we return here
                    return;
                } else {
                    throw popupError;
                }
            }

            const firebaseUser = result.user;

            // Log successful authentication details
            console.log("Google sign-in successful:", {
                uid: firebaseUser.uid,
                email: firebaseUser.email,
                displayName: firebaseUser.displayName,
                photoURL: firebaseUser.photoURL,
            });

            showSuccess("Accesso con Google completato!");

            // Clear any existing form data
            clearForm();

            createUserDocument(firebaseUser.uid);
            setTimeout(() => {
                push(redirectUrl);
            }, 1500);
        } catch (error: any) {
            console.error("Google sign in error:", error);

            // More specific error handling for Google sign-in
            if (error.code === "auth/popup-closed-by-user") {
                showError("Accesso annullato dall'utente");
            } else if (error.code === "auth/popup-blocked") {
                showError(
                    "Popup bloccato dal browser. Abilita i popup e riprova",
                );
            } else if (error.code === "auth/network-request-failed") {
                showError(
                    "Errore di connessione. Controlla la tua connessione internet",
                );
            } else if (error.code === "auth/internal-error") {
                showError("Errore interno. Riprova più tardi");
            } else {
                showError(handleFirebaseError(error as AuthError));
            }
        } finally {
            isLoading = false;
        }
    }

    // Handle Enter key press
    function handleKeyPress(event: KeyboardEvent) {
        if (event.key === "Enter") {
            if (showForgotPassword) {
                handlePasswordReset();
            } else {
                handleSubmit();
            }
        }
    }

    onMount(() => {
        // Parse query parameters from the current location
        const params = new URLSearchParams($querystring);
        const redirect = params.get("redirect");

        if (redirect) {
            redirectUrl = decodeURIComponent(redirect);
        }
        // Initialize auth store
        initAuthStore();

        // Check for redirect result (for mobile/popup-blocked scenarios)
        getRedirectResult(auth)
            .then((result) => {
                if (result) {
                    // User successfully signed in via redirect
                    const firebaseUser = result.user;
                    console.log("Redirect sign-in successful:", firebaseUser);
                    showSuccess("Accesso completato!");

                    setTimeout(() => {
                        push(redirectUrl);
                    }, 1500);
                }
            })
            .catch((error) => {
                console.error("Redirect result error:", error);
                if (error.code !== "auth/popup-closed-by-user") {
                    showError(handleFirebaseError(error as AuthError));
                }
            });

        // Listen for auth state changes
        const unsubscribe = onAuthStateChanged(auth, (firebaseUser) => {
            if (firebaseUser) {
                // Check if email is verified (skip for Google/Facebook sign-ins)
                if (
                    firebaseUser.emailVerified ||
                    firebaseUser.providerData.some(
                        (provider) =>
                            provider.providerId === "google.com" ||
                            provider.providerId === "facebook.com",
                    )
                ) {
                    // The auth store will handle the redirect automatically
                    showSuccess(
                        "Autenticazione completata! Reindirizzamento...",
                    );
                    setTimeout(() => {
                        push(redirectUrl);
                    }, 1500);
                } else {
                    // Show email verification message only for email/password accounts
                    if (!showForgotPassword) {
                        showEmailVerificationMessage = true;
                    }
                }
            } else {
                showEmailVerificationMessage = false;
            }
        });

        // Cleanup listener on component destroy
        return () => unsubscribe();
    });
</script>

<!-- The HTML template -->
<div
    class="min-h-screen bg-gradient-to-b from-purple-400 to-[#1e1f25] flex items-center justify-center p-2 sm:p-6 overflow-x-hidden"
>
    <!-- Back to Home Button -->
    <button
        on:click={() => push("/")}
        class="fixed top-4 left-4 bg-white/20 hover:bg-white/30 text-white px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg z-50 transition-all duration-300 font-medium text-sm sm:text-base backdrop-blur-sm"
    >
        ←
    </button>

    <!-- Error Popup -->
    {#if showErrorPopup}
        <div
            class="fixed top-4 left-1/2 transform -translate-x-1/2 bg-red-500 text-white px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg z-50 transition-all duration-300 max-w-[90vw]"
            transition:slide={{ duration: 300 }}
        >
            <div class="flex items-center justify-between gap-2 sm:gap-4">
                <span class="text-sm sm:text-base">{errorMessage}</span>
                <button
                    on:click={() => (showErrorPopup = false)}
                    class="text-white hover:text-red-200 font-bold text-lg leading-none min-w-[24px]"
                >
                    ×
                </button>
            </div>
        </div>
    {/if}

    <!-- Success Popup -->
    {#if showSuccessPopup}
        <div
            class="fixed top-4 left-1/2 transform -translate-x-1/2 bg-green-500 text-white px-4 py-2 sm:px-6 sm:py-3 rounded-2xl shadow-lg z-50 transition-all duration-300 max-w-[90vw]"
            transition:slide={{ duration: 300 }}
        >
            <div class="flex items-center justify-between gap-2 sm:gap-4">
                <span class="text-sm sm:text-base">{successMessage}</span>
                <button
                    on:click={() => (showSuccessPopup = false)}
                    class="text-white hover:text-green-200 font-bold text-lg leading-none min-w-[24px]"
                >
                    ×
                </button>
            </div>
        </div>
    {/if}

    <!-- Main Container -->
    <div class="w-full max-w-md mx-auto">
        <div
            class="bg-[#1e1f25] text-white rounded-2xl shadow-lg p-6 sm:p-8 transition-all duration-500"
            in:fly={{ x: 1000, duration: 500, delay: 100 }}
        >
            <!-- Header -->
            <div class="text-center mb-8">
                <h1 class="text-3xl sm:text-4xl font-bold leading-tight mb-2">
                    {#if showForgotPassword}
                        Reset Password 🔑
                    {:else if isSignUp}
                        Benvenuto! 🎉
                    {:else}
                        Bentornato! 👋
                    {/if}
                </h1>
                <p class="text-gray-400 text-sm">
                    {#if showForgotPassword}
                        Inserisci la tua email per resettare la password
                    {:else if isSignUp}
                        Crea il tuo account per iniziare
                    {:else}
                        Accedi al tuo account
                    {/if}
                </p>
            </div>

            {#if showEmailVerificationMessage}
                <div
                    class="bg-blue-500/20 border border-blue-500/50 rounded-lg p-4 mb-6"
                    transition:slide={{ duration: 300 }}
                >
                    <div class="text-center">
                        <p class="text-blue-300 text-sm mb-4">
                            📧 Controlla la tua email per verificare l'account
                        </p>
                        <div
                            class="flex flex-col sm:flex-row gap-3 justify-center"
                        >
                            <button
                                on:click={resendEmailVerification}
                                disabled={isLoading}
                                class="bg-purple-500/20 hover:bg-purple-500/30 border border-purple-400/50 text-purple-300 hover:text-purple-200 px-4 py-2 rounded-lg text-sm font-medium disabled:opacity-50 transition-all duration-200 hover:scale-105"
                            >
                                🔄 Invia nuovamente email
                            </button>
                            <button
                                on:click={checkEmailVerification}
                                disabled={isLoading}
                                class="bg-green-500/20 hover:bg-green-500/30 border border-green-400/50 text-green-300 hover:text-green-200 px-4 py-2 rounded-lg text-sm font-medium disabled:opacity-50 transition-all duration-200 hover:scale-105"
                            >
                                ✅ Ho verificato l'email
                            </button>
                        </div>
                    </div>
                </div>
            {/if}

            <!-- Form -->
            {#if showForgotPassword}
                <!-- Password Reset Form -->
                <form
                    on:submit|preventDefault={handlePasswordReset}
                    class="space-y-4"
                >
                    <div class="flex flex-col">
                        <label
                            class="text-sm font-semibold mb-2"
                            for="resetEmail"
                        >
                            Email:
                        </label>
                        <input
                            id="resetEmail"
                            type="email"
                            class="w-full border rounded-lg px-4 py-3 text-white bg-transparent text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all duration-200"
                            class:border-white={resetEmailValid}
                            class:border-red-500={!resetEmailValid}
                            bind:value={resetEmail}
                            required
                            on:keypress={handleKeyPress}
                            on:blur={() =>
                                (resetEmailValid = validateEmail(resetEmail))}
                        />
                        {#if !resetEmailValid}
                            <span class="text-red-400 text-xs mt-1"
                                >Inserisci un'email valida</span
                            >
                        {/if}
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        class="w-full bg-purple-500 hover:bg-purple-600 disabled:bg-gray-500 disabled:cursor-not-allowed text-white font-semibold py-3 px-6 rounded-xl transition-all duration-200 transform hover:scale-105 active:scale-95 mt-6"
                    >
                        {#if isLoading}
                            <div
                                class="flex items-center justify-center space-x-2"
                            >
                                <div
                                    class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"
                                ></div>
                                <span>Caricamento...</span>
                            </div>
                        {:else}
                            Invia email di reset
                        {/if}
                    </button>
                </form>
            {:else}
                <!-- Login/Register Form -->
                <form on:submit|preventDefault={handleSubmit} class="space-y-4">
                    <!-- Name fields (only for sign up) -->
                    {#if isSignUp}
                        <div
                            class="grid grid-cols-2 gap-4"
                            transition:slide={{ duration: 300 }}
                        >
                            <div class="flex flex-col">
                                <label
                                    class="text-sm font-semibold mb-2"
                                    for="firstName"
                                >
                                    Nome:
                                </label>
                                <input
                                    id="firstName"
                                    type="text"
                                    class="w-full border border-white rounded-lg px-4 py-3 text-white bg-transparent text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all duration-200"
                                    bind:value={firstName}
                                    required
                                    on:keypress={handleKeyPress}
                                />
                            </div>
                            <div class="flex flex-col">
                                <label
                                    class="text-sm font-semibold mb-2"
                                    for="lastName"
                                >
                                    Cognome:
                                </label>
                                <input
                                    id="lastName"
                                    type="text"
                                    class="w-full border border-white rounded-lg px-4 py-3 text-white bg-transparent text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all duration-200"
                                    bind:value={lastName}
                                    required
                                    on:keypress={handleKeyPress}
                                />
                            </div>
                        </div>
                    {/if}

                    <!-- Email -->
                    <div class="flex flex-col">
                        <label class="text-sm font-semibold mb-2" for="email">
                            Email:
                        </label>
                        <input
                            id="email"
                            type="email"
                            class="w-full border rounded-lg px-4 py-3 text-white bg-transparent text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all duration-200"
                            class:border-white={emailValid}
                            class:border-red-500={!emailValid}
                            bind:value={email}
                            required
                            on:keypress={handleKeyPress}
                            on:blur={() => (emailValid = validateEmail(email))}
                        />
                        {#if !emailValid}
                            <span class="text-red-400 text-xs mt-1"
                                >Inserisci un'email valida</span
                            >
                        {/if}
                    </div>

                    <!-- Password -->
                    <div class="flex flex-col">
                        <label
                            class="text-sm font-semibold mb-2"
                            for="password"
                        >
                            Password:
                        </label>
                        <input
                            id="password"
                            type="password"
                            class="w-full border rounded-lg px-4 py-3 text-white bg-transparent text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all duration-200"
                            class:border-white={passwordValid}
                            class:border-red-500={!passwordValid}
                            bind:value={password}
                            required
                            on:keypress={handleKeyPress}
                            on:blur={() =>
                                (passwordValid = validatePassword(password))}
                        />
                        {#if !passwordValid}
                            <span class="text-red-400 text-xs mt-1"
                                >La password deve contenere almeno 8 caratteri</span
                            >
                        {/if}
                    </div>

                    <!-- Confirm Password (only for sign up) -->
                    {#if isSignUp}
                        <div
                            class="flex flex-col"
                            transition:slide={{ duration: 300 }}
                        >
                            <label
                                class="text-sm font-semibold mb-2"
                                for="confirmPassword"
                            >
                                Conferma Password:
                            </label>
                            <input
                                id="confirmPassword"
                                type="password"
                                class="w-full border rounded-lg px-4 py-3 text-white bg-transparent text-sm focus:outline-none focus:ring-2 focus:ring-purple-400 focus:border-transparent transition-all duration-200"
                                class:border-white={confirmPasswordValid}
                                class:border-red-500={!confirmPasswordValid}
                                bind:value={confirmPassword}
                                required
                                on:keypress={handleKeyPress}
                                on:blur={() =>
                                    (confirmPasswordValid =
                                        password === confirmPassword)}
                            />
                            {#if !confirmPasswordValid}
                                <span class="text-red-400 text-xs mt-1"
                                    >Le password non coincidono</span
                                >
                            {/if}
                        </div>
                    {/if}

                    <!-- Terms and Conditions (only for sign up) -->
                    {#if isSignUp}
                        <div
                            class="flex items-start space-x-3"
                            transition:slide={{ duration: 300 }}
                        >
                            <input
                                id="acceptTerms"
                                type="checkbox"
                                class="mt-1 w-4 h-4 accent-purple-400"
                                bind:checked={acceptTerms}
                            />
                            <label
                                for="acceptTerms"
                                class="text-sm text-gray-300"
                            >
                                Accetto i <a
                                    href="#"
                                    class="text-purple-400 hover:text-purple-300 underline"
                                    >termini e condizioni</a
                                >
                                e la
                                <a
                                    href="#"
                                    class="text-purple-400 hover:text-purple-300 underline"
                                    >privacy policy</a
                                >
                            </label>
                        </div>
                    {/if}

                    <!-- Submit Button -->
                    <button
                        type="submit"
                        disabled={isLoading}
                        class="w-full bg-purple-500 hover:bg-purple-600 disabled:bg-gray-500 disabled:cursor-not-allowed text-white font-semibold py-3 px-6 rounded-xl transition-all duration-200 transform hover:scale-105 active:scale-95 mt-6"
                    >
                        {#if isLoading}
                            <div
                                class="flex items-center justify-center space-x-2"
                            >
                                <div
                                    class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"
                                ></div>
                                <span>Caricamento...</span>
                            </div>
                        {:else}
                            {isSignUp ? "Registrati" : "Accedi"}
                        {/if}
                    </button>
                </form>
            {/if}

            <!-- Toggle Mode / Forgot Password -->
            <div class="text-center mt-6">
                {#if showForgotPassword}
                    <p class="text-gray-400 text-sm">
                        Ricordi la password?
                        <button
                            on:click={toggleForgotPassword}
                            class="text-purple-400 hover:text-purple-300 font-semibold ml-1 underline transition-colors duration-200"
                        >
                            Torna al login
                        </button>
                    </p>
                {:else}
                    <p class="text-gray-400 text-sm">
                        {#if isSignUp}
                            Hai già un account?
                        {:else}
                            Non hai un account?
                        {/if}
                        <button
                            on:click={toggleMode}
                            class="text-purple-400 hover:text-purple-300 font-semibold ml-1 underline transition-colors duration-200"
                        >
                            {isSignUp ? "Accedi" : "Registrati"}
                        </button>
                    </p>

                    <!-- Forgot Password Link (only for login) -->
                    {#if !isSignUp}
                        <p class="text-gray-400 text-sm mt-2">
                            <button
                                on:click={toggleForgotPassword}
                                class="text-purple-400 hover:text-purple-300 font-semibold underline transition-colors duration-200"
                            >
                                Password dimenticata?
                            </button>
                        </p>
                    {/if}
                {/if}
            </div>

            <!-- Divider and Social Login (only for login/register, not password reset) -->
            {#if !showForgotPassword}
                <!-- Divider -->
                <div class="flex items-center my-6">
                    <div class="flex-1 border-t border-gray-600"></div>
                    <span class="px-4 text-gray-400 text-sm">oppure</span>
                    <div class="flex-1 border-t border-gray-600"></div>
                </div>

                <!-- Social Login -->
                <div class="space-y-3">
                    <button
                        type="button"
                        on:click={handleGoogleLogin}
                        disabled={isLoading}
                        class="w-full bg-white hover:bg-gray-100 disabled:bg-gray-300 disabled:cursor-not-allowed text-gray-900 font-semibold py-3 px-6 rounded-xl transition-all duration-200 transform hover:scale-105 active:scale-95 flex items-center justify-center space-x-2"
                    >
                        <svg class="w-5 h-5" viewBox="0 0 24 24">
                            <path
                                fill="currentColor"
                                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                            />
                            <path
                                fill="currentColor"
                                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                            />
                            <path
                                fill="currentColor"
                                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                            />
                            <path
                                fill="currentColor"
                                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                            />
                        </svg>
                        <span>Continua con Google</span>
                    </button>
                </div>
            {/if}
        </div>
    </div>
</div>

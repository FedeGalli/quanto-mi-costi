<script lang="ts">
    export let isLoadingSaving = false;
    export let house_name = "";
    export let showNamePopup = false;
    export let handleSaveHouse: (arg0: string) => Promise<boolean>;
    export let handleOverwriteHouse: (arg0: string) => Promise<boolean>;
    let currentHouseName = house_name;

    let errorMessage = "";
    let showOverwriteConfirm = false;

    async function onSaveClick() {
        errorMessage = ""; // Clear previous errors
        showOverwriteConfirm = false;

        const result = await handleSaveHouse(currentHouseName);

        if (result === false) {
            showOverwriteConfirm = true;
            errorMessage = "Esiste già una casa salvata con questo nome.";
        }
    }

    async function confirmOverwrite() {
        showOverwriteConfirm = false;
        errorMessage = "";

        const result = await handleOverwriteHouse(currentHouseName);

        if (result) {
            closePopup();
        } else {
            errorMessage = "Errore durante la sovrascrittura. Riprova.";
        }
    }

    function cancelOverwrite() {
        showOverwriteConfirm = false;
        errorMessage = "";
    }

    function handleKeydown(event: any) {
        if (event.key === "Enter") {
            onSaveClick();
        } else if (event.key === "Escape") {
            showNamePopup = false;
            currentHouseName = "";
            errorMessage = "";
            showOverwriteConfirm = false;
        }
    }

    function closePopup() {
        showNamePopup = false;
        currentHouseName = "";
        errorMessage = "";
        showOverwriteConfirm = false;
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    on:click={closePopup}
>
    <div
        class="bg-white rounded-2xl p-6 w-full max-w-sm shadow-2xl transform transition-all duration-300 scale-100"
        on:click|stopPropagation
    >
        <!-- Icon and Title -->
        <div class="text-center mb-6">
            <div
                class="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4"
            >
                <svg
                    class="w-8 h-8 text-purple-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"
                    />
                    <polyline points="9,22 9,12 15,12 15,22" />
                </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">
                Salva la tua casa
            </h3>
            <p class="text-gray-600 text-sm">
                Inserisci un nome per identificare questa configurazione
            </p>
        </div>

        <!-- Input with simple label -->
        <div class="mb-6">
            <label
                for="house-name-input"
                class="block text-sm font-medium text-gray-700 mb-2"
            >
                Nome della casa
            </label>
            <input
                id="house-name-input"
                type="text"
                bind:value={currentHouseName}
                on:keydown={handleKeydown}
                on:input={() => {
                    errorMessage = "";
                    showOverwriteConfirm = false;
                }}
                placeholder={house_name == ""
                    ? "Es. Casa al mare, Villa in montagna..."
                    : ""}
                class="w-full px-4 py-3 border-2 rounded-xl focus:ring-0 outline-none transition-all duration-200 text-gray-800 placeholder-gray-400
                    {errorMessage
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-200 focus:border-purple-500'}"
                disabled={isLoadingSaving}
            />

            <!-- Error Message with Overwrite Options -->
            {#if errorMessage}
                <div class="mt-2">
                    <div
                        class="flex items-center gap-2 text-red-600 text-sm mb-3"
                    >
                        <svg
                            class="w-4 h-4 flex-shrink-0"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                                clip-rule="evenodd"
                            />
                        </svg>
                        <span>{errorMessage}</span>
                    </div>

                    {#if showOverwriteConfirm}
                        <div
                            class="bg-yellow-50 border border-yellow-200 rounded-lg p-3"
                        >
                            <p class="text-yellow-800 text-sm mb-3">
                                Vuoi sovrascrivere la casa esistente con questi
                                nuovi valori?
                            </p>
                            <div class="flex gap-2">
                                <button
                                    on:click={confirmOverwrite}
                                    class="px-3 py-1 bg-yellow-600 hover:bg-yellow-700 text-white text-sm rounded transition-colors duration-200"
                                    disabled={isLoadingSaving}
                                >
                                    Sovrascrivi
                                </button>
                                <button
                                    on:click={cancelOverwrite}
                                    class="px-3 py-1 bg-gray-300 hover:bg-gray-400 text-gray-700 text-sm rounded transition-colors duration-200"
                                    disabled={isLoadingSaving}
                                >
                                    Annulla
                                </button>
                            </div>
                        </div>
                    {/if}
                </div>
            {/if}
        </div>

        <!-- Action Buttons -->
        <div class="flex space-x-3">
            <button
                on:click={closePopup}
                class="flex-1 px-4 py-3 text-gray-600 border-2 border-gray-200 rounded-xl hover:bg-gray-50 transition-all duration-200 font-medium"
                disabled={isLoadingSaving}
            >
                Annulla
            </button>
            <button
                on:click={onSaveClick}
                class="flex-1 px-4 py-3 bg-gradient-to-r from-purple-600 to-purple-700 hover:from-purple-700 hover:to-purple-800 text-white rounded-xl transition-all duration-200 font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center shadow-lg shadow-purple-500/25"
                disabled={isLoadingSaving || !currentHouseName.trim()}
            >
                {#if isLoadingSaving}
                    <svg
                        class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
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
                        />
                        <path
                            class="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        />
                    </svg>
                    Salvando...
                {:else}
                    <svg
                        class="w-4 h-4 mr-2"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 13l4 4L19 7"
                        />
                    </svg>
                    Salva
                {/if}
            </button>
        </div>
    </div>
</div>

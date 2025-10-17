<script lang="ts">
    import { _ } from "svelte-i18n";

    export let isDeleting = false;
    export let showDeletePopup = false;
    export let handleDeleteUser: () => Promise<boolean>;
    let passKey = "";

    let errorMessage = "";

    async function onDeleteClick() {
        errorMessage = ""; // Clear previous errors

        if (passKey != "Delete") {
            errorMessage = $_("error.userDelete.typeError");
            return;
        }

        const result = await handleDeleteUser();

        if (result === false) {
            errorMessage = $_("error.userDelete.generalError");
        }
    }

    function handleKeydown(event: any) {
        if (event.key === "Enter") {
            onDeleteClick();
        } else if (event.key === "Escape") {
            showDeletePopup = false;
            passKey = "";
            errorMessage = "";
        }
    }

    function closePopup() {
        showDeletePopup = false;
        passKey = "";
        errorMessage = "";
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
                class="w-16 h-16 bg-red-200 rounded-full flex items-center justify-center mx-auto mb-4"
            >
                <svg
                    class="w-8 h-8 text-red-600"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M3 6h18" />
                    <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                    <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                    <line x1="10" y1="11" x2="10" y2="17" />
                    <line x1="14" y1="11" x2="14" y2="17" />
                </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">
                {$_("userDelete.delete")}
            </h3>
            <p class="text-gray-600 text-sm">
                {$_("userDelete.deleteInfo")}
            </p>
        </div>

        <!-- Input with simple label -->
        <div class="mb-6">
            <label
                for="house-name-input"
                class="block text-sm font-medium text-gray-700 mb-2"
            >
                {$_("userDelete.type")}
            </label>
            <input
                id="house-name-input"
                type="text"
                bind:value={passKey}
                on:keydown={handleKeydown}
                on:input={() => {
                    errorMessage = "";
                }}
                class="w-full px-4 py-3 border-2 rounded-xl focus:ring-0 outline-none transition-all duration-200 text-gray-800 placeholder-gray-400
                    {errorMessage
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-200 focus:border-purple-500'}"
                disabled={isDeleting}
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
                </div>
            {/if}
        </div>

        <!-- Action Buttons -->
        <div class="flex space-x-3">
            <button
                on:click={closePopup}
                class="flex-1 px-4 py-3 text-gray-600 border-2 border-gray-200 rounded-xl hover:bg-gray-50 transition-all duration-200 font-medium"
                disabled={isDeleting}
            >
                {$_("houseSaving.abort")}
            </button>
            <button
                on:click={onDeleteClick}
                class="flex-1 px-4 py-3 bg-gradient-to-r from-red-600 to-red-800 hover:from-red-600 hover:to-red-800 text-white rounded-xl transition-all duration-200 font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center shadow-lg shadow-purple-500/25"
                disabled={isDeleting || !passKey.trim()}
            >
                {#if isDeleting}
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
                    {$_("userDelete.deleting")}
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
                    {$_("userDelete.deleteButton")}
                {/if}
            </button>
        </div>
    </div>
</div>

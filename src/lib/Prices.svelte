<script lang="ts">
    import { fade, slide } from "svelte/transition";
    import { onMount } from "svelte";

    let selectedComune: string = "";
    let selectedZona: string = "";
    let selectedTipologia: string = "";
    let selectedStato: string = "";
    let selectedMq: number = 0;

    let comuniOptions: any[] = [];
    let comuniInfo: any[] = [];
    let zoneOptions: any[] = [];
    let tipologieOptions: any[] = [];
    let statoOptions: any[] = [];
    let testData: any[] = [];

    $: if (selectedZona) {
        selectZona(selectedZona);
    }
    $: if (selectedTipologia) {
        selectTipologia(selectedTipologia);
    }

    // Dropdown visibility state
    let showComuniDropdown: boolean = false;

    function getMunicipalitiesString(): string {
        return "http://localhost:8000/get-municipalities-list";
    }
    function getMunicipalitiesInfoString(com: string): string {
        return "http://127.0.0.1:8000/get-municipalities-info?com=" + com;
    }
    function getPriceVolumeDataString(
        com: string,
        zona: string,
        tipo: string,
        stato: string,
        mq: number,
    ): string {
        return (
            "http://127.0.0.1:8000/get-price-volume-data?com=" +
            com +
            "&zone=" +
            zona +
            "&type=" +
            tipo +
            "&state=" +
            stato +
            "&mq=" +
            mq
        );
    }

    async function getMunicipalitiesData() {
        const apiStringUrl: string = getMunicipalitiesString();
        let responseCode: string = "200";

        try {
            const response = await fetch(apiStringUrl);

            if (!response.ok || response.status == 429) {
                throw new Error(`${response.status}`);
            }
            const data = await response.json();
            comuniOptions = data["DATA"];
        } catch (error: any) {
            responseCode = error.message;
        }
    }

    async function getPriceVolumeData(
        com: string,
        zona: string,
        tipo: string,
        stato: string,
        mq: number,
    ) {
        const apiStringUrl: string = getPriceVolumeDataString(
            com,
            zona,
            tipo,
            stato,
            mq,
        );
        let responseCode: string = "200";

        try {
            const response = await fetch(apiStringUrl);

            if (!response.ok || response.status == 429) {
                throw new Error(`${response.status}`);
            }
            const data = await response.json();
            console.log(data);
            testData = data;
        } catch (error: any) {
            responseCode = error.message;
        }
    }

    async function getMunicipalitiesInfoData(com: string) {
        const apiStringUrl: string = getMunicipalitiesInfoString(com);
        let responseCode: string = "200";

        try {
            const response = await fetch(apiStringUrl);

            if (!response.ok || response.status == 429) {
                throw new Error(`${response.status}`);
            }
            const data = await response.json();
            comuniInfo = data["DATA"];
        } catch (error: any) {
            responseCode = error.message;
        }
    }

    function getZoneOptions() {
        zoneOptions = [];
        const zoneSet = new Set<string>(); // or Set<any> if not string
        comuniInfo.forEach((prop) => {
            zoneSet.add(prop[1]);
        });
        zoneOptions = Array.from(zoneSet);
    }

    function getTypeOptions(zona: string) {
        tipologieOptions = [];
        const tipologieSet = new Set<string>(); // or Set<any> if not string
        comuniInfo.forEach((prop) => {
            if (prop[1] == zona) {
                tipologieSet.add(prop[2]);
            }
            console.log(prop[2]);
        });
        tipologieOptions = Array.from(tipologieSet);
    }

    function getStatoOptions(zona: string) {
        statoOptions = [];
        const statoSet = new Set<string>(); // or Set<any> if not string
        comuniInfo.forEach((prop) => {
            if (prop[2] == zona) {
                statoSet.add(prop[3]);
            }
        });
        statoOptions = Array.from(statoSet);
    }

    let filteredComuni: string[] = [];

    function filterComuni(): void {
        if (selectedComune.length === 0) {
            filteredComuni = comuniOptions.slice(0, 10); // Show first 10 by default
        } else {
            filteredComuni = comuniOptions
                .filter((comune: string) =>
                    comune.toLowerCase().includes(selectedComune.toLowerCase()),
                )
                .slice(0, 10); // Limit to 10 results for performance
        }
    }

    function comuneChange() {
        selectedZona = "";
        selectedTipologia = "";
        selectedStato = "";
    }

    // Selection functions
    async function selectComune(comune: string) {
        selectedComune = comune;
        showComuniDropdown = false;
        comuneChange();
        await getMunicipalitiesInfoData(comune);
        getZoneOptions();
    }

    function zonaChange() {
        selectedTipologia = "";
        selectedStato = "";
    }

    function selectZona(zona: string): void {
        zonaChange();
        getTypeOptions(zona);
    }

    function tipologiaChange() {
        selectedStato = "";
    }

    function selectTipologia(zona: string): void {
        tipologiaChange();
        getStatoOptions(zona);
    }

    onMount(async () => {
        await getMunicipalitiesData();
        filterComuni();
    });
</script>

<!-- Prices section -->
<div transition:fade={{ duration: 500 }}>
    <!-- Input section -->
    <div
        class="flex flex-col items-center gap-4 mb-4"
        transition:slide={{ duration: 500 }}
    >
        <!-- Question -->
        <h2
            class="font-bold text-white leading-tight text-base sm:text-lg text-center mb-2"
            transition:slide={{
                duration: 500,
            }}
        >
            Quanto <span class="font-bold text-purple-400">costano</span>
            le case nella mia
            <span class="font-bold text-purple-400">zona</span>?
        </h2>
        <!-- First row: Main property inputs -->
        <div class="flex justify-center items-center gap-4 flex-wrap">
            <!-- Comune (Autocomplete) -->
            <div class="flex flex-col items-center gap-2">
                <h3
                    class="font-bold text-white text-sm sm:text-base text-center"
                >
                    Comune
                </h3>
                <div class="relative">
                    <input
                        type="text"
                        class="w-40 border border-white rounded px-3 py-2 text-white bg-transparent focus:border-white focus:outline-none placeholder-gray-300"
                        placeholder="Cerca comune..."
                        bind:value={selectedComune}
                        on:input={filterComuni}
                        on:focus={() => (showComuniDropdown = true)}
                        on:blur={() =>
                            setTimeout(() => (showComuniDropdown = false), 200)}
                    />
                    {#if showComuniDropdown && filteredComuni.length > 0}
                        <div
                            class="absolute top-full left-0 right-0 bg-white border border-gray-300 rounded-b max-h-40 overflow-y-auto z-10"
                        >
                            {#each filteredComuni as comune}
                                <div
                                    class="px-3 py-2 hover:bg-gray-100 cursor-pointer text-black"
                                    on:click={() => selectComune(comune)}
                                >
                                    {comune}
                                </div>
                            {/each}
                        </div>
                    {/if}
                </div>
            </div>

            <!-- Zona (Dropdown) -->
            <div class="flex flex-col items-center gap-2">
                <h3
                    class="font-bold text-white text-sm sm:text-base text-center"
                >
                    Zona
                </h3>
                <div class="relative">
                    <select
                        class="w-32 border border-white rounded px-3 py-2 text-white bg-transparent focus:border-white focus:outline-none appearance-none cursor-pointer"
                        bind:value={selectedZona}
                    >
                        <option value="" class="text-black">Seleziona</option>
                        {#each zoneOptions as zona}
                            <option value={zona} class="text-black"
                                >{zona}</option
                            >
                        {/each}
                    </select>
                    <div class="absolute right-3 top-2 pointer-events-none">
                        <svg
                            class="w-4 h-4 text-white"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 9l-7 7-7-7"
                            ></path>
                        </svg>
                    </div>
                </div>
            </div>

            <!-- Tipologia (Dropdown) -->
            <div class="flex flex-col items-center gap-2">
                <h3
                    class="font-bold text-white text-sm sm:text-base text-center"
                >
                    Tipologia
                </h3>
                <div class="relative">
                    <select
                        class="w-32 border border-white rounded px-3 py-2 text-white bg-transparent focus:border-white focus:outline-none appearance-none cursor-pointer"
                        bind:value={selectedTipologia}
                    >
                        <option value="" class="text-black">Seleziona</option>
                        {#each tipologieOptions as tipologia}
                            <option value={tipologia} class="text-black"
                                >{tipologia}</option
                            >
                        {/each}
                    </select>
                    <div class="absolute right-3 top-2 pointer-events-none">
                        <svg
                            class="w-4 h-4 text-white"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 9l-7 7-7-7"
                            ></path>
                        </svg>
                    </div>
                </div>
            </div>

            <!-- Stato (Dropdown) -->
            <div class="flex flex-col items-center gap-2">
                <h3
                    class="font-bold text-white text-sm sm:text-base text-center"
                >
                    Stato
                </h3>
                <div class="relative">
                    <select
                        class="w-32 border border-white rounded px-3 py-2 text-white bg-transparent focus:border-white focus:outline-none appearance-none cursor-pointer"
                        bind:value={selectedStato}
                    >
                        <option value="" class="text-black">Seleziona</option>
                        {#each statoOptions as stato}
                            <option value={stato} class="text-black"
                                >{stato}</option
                            >
                        {/each}
                    </select>
                    <div class="absolute right-3 top-2 pointer-events-none">
                        <svg
                            class="w-4 h-4 text-white"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 9l-7 7-7-7"
                            ></path>
                        </svg>
                    </div>
                </div>
            </div>

            <!-- m² (Number input) -->
            <div class="flex flex-col items-center gap-2">
                <h3
                    class="font-bold text-white text-sm sm:text-base text-center"
                >
                    m²
                </h3>
                <div class="relative">
                    <input
                        type="number"
                        class="w-24 border border-white rounded px-3 py-2 text-white bg-transparent focus:border-white focus:outline-none"
                        min="0"
                        step="5"
                        max="500"
                        bind:value={selectedMq}
                    />
                </div>
            </div>
            <button
                on:click={() => {
                    getPriceVolumeData(
                        selectedComune,
                        selectedZona,
                        selectedTipologia,
                        selectedStato,
                        selectedMq,
                    );
                }}
            >
                CALCULATE
            </button>
        </div>
    </div>
</div>

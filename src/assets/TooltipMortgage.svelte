<script lang="ts">
    export let showTooltip: boolean[] = [];
    export let index: number = 0;
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
            Informazioni principali per un <span class="text-purple-400"
                >mutuo</span
            >...
        </h2>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1.</span> Importo mutuo:
            <span class="text-base font-medium leading-tight"
                >La totalità del capitale richiesto alla banca.</span
            >
        </p>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">2.</span> Durata (Anni):
            <span class="text-base font-medium leading-tight">
                Il tempo necessario con il quale si desidera ripagare il
                capitale richiesto.</span
            >
        </p>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">3.</span> TAEG:
            <span class="text-base font-medium leading-tight">
                Il tasso di interesse annuale applicato al capitale residuo da
                restituire. Il TAEG a differenza del TAN è comprensivo di tutti
                i costi reali che vengono applicati, anche quelli assicurativi
                per polizza scoppio inecendio obbligatoria / polize vita.</span
            >
        </p>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">4.</span> Per aumentare le chache di ottenere
            un mutuo:
        </p>

        <ul class="list-disc pl-6 text-base leading-tight space-y-3">
            <li>
                <span class="font-bold leading-tight"
                    >Reddito e capacità di rimborso</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        <span class="font-bold">Stabilità lavorativa:</span> Lavoro
                        a tempo indeterminato, da almeno 6–12 mesi. In caso di lavoro
                        autonomo o partita IVA, servono almeno 2 anni di attività
                        con redditi dichiarati stabili.
                    </li>
                    <li>
                        <span class="font-bold">Rapporto rata/reddito:</span> La
                        rata mensile del mutuo non deve superare il 30–35% del reddito
                        netto mensile familiare.
                    </li>
                    <li>
                        <span class="font-bold"
                            >Reddito familiare aggregato:</span
                        > Se ci sono più intestatari del mutuo, si considera il reddito
                        complessivo.
                    </li>
                </ul>
            </li>
            <li>
                <span class="font-bold leading-tight">Età del richiedente</span>
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        In genere, il mutuo deve terminare prima dei 75–80 anni
                        di età del richiedente.
                    </li>
                </ul>
            </li>
            <li>
                <span class="font-bold leading-tight"
                    >Importo richiesto e valore dell’immobile</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        <span class="font-bold"
                            >Rapporto LTV (Loan-to-Value):</span
                        > L’importo del mutuo viene calcolato rispetto al valore
                        dell’immobile. Di solito le banche finanziano fino all’80%
                        del valore di perizia.
                    </li>
                    <li>
                        <span class="font-bold">Cataso e ipoteche:</span> L'immobile
                        deve essere regolarmente accatastato, libero da ipoteche
                        pregresse e in buono stato.
                    </li>
                </ul>
            </li>
            <li>
                <span class="font-bold leading-tight"
                    >Affidabilità creditizia</span
                >
                <ul
                    class="list-disc pl-6 mt-2 text-base leading-tight space-y-1"
                >
                    <li>
                        <span class="font-bold"
                            >Contorollo situazione creditizia:</span
                        > Le banche consultano dei database per verificare se hai
                        altri prestiti in corso ed eventuali ritardi nei pagamenti.
                    </li>
                    <li>
                        <span class="font-bold">Credit history:</span>
                        Avere una cronologia pregressa di pagamenti regolari o mutui/prestiti
                        completamente ripagati senza intoppi, aumenta la tua credibilità
                        e affidabilità agli occhi della banca.
                    </li>
                </ul>
            </li>
        </ul>
    </div>
</div>

<script lang="ts">
    export let showTooltip: boolean = false;
</script>

<!-- Tooltip (Overlay) -->
<div
    class="absolute top-0 left-0 w-full h-full bg-[#1e1f25] rounded-2xl shadow-xl p-8 overflow-auto z-20 transition-all transform duration-600 ease-out-[cubic-bezier(0.22, 1, 0.36, 1)]"
    class:opacity-100={showTooltip}
    class:opacity-0={!showTooltip}
    class:-translate-y-500={!showTooltip}
    class:translate-y-0={showTooltip}
>
    <div class="space-y-6">
        <h2 class="text-4xl font-bold leading-tight">
            Hai l'agevolazione <span class="text-purple-400">prima casa</span>
            se...
        </h2>
        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">1.</span> Requisiti sull'immobile acquistato
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                Non appartenere alle categorie catastali di lusso: quindi NO A1
                (abitazioni signorili), A8 (ville), A9 (castelli).
            </li>
            <li>
                Essere ubicato nel comune dove hai o avrai la residenza (entro
                18 mesi dall'acquisto).
            </li>
            <li>
                Se vendi la casa acquistata con agevolazioni entro 5 anni, per
                non perdere il beneficio devi riacquistare entro 1 anno un'altra
                prima casa.
            </li>
        </ul>

        <p class="text-xl font-bold leading-tight">
            <span class="text-purple-400">2.</span> Requisiti personali dell’acquirente
        </p>

        <ul class="list-disc pl-6 text-base font-bold leading-tight space-y-2">
            <li>
                Non deve essere titolare, esclusivo o in comunione con il
                coniuge, di altra abitazione nel territorio del comune dove si
                trova l’immobile.
            </li>
            <li>
                Non deve possedere altri immobili acquistati con le agevolazioni
                prima casa, su tutto il territorio nazionale, a meno che li
                venda entro 12 mesi dall’acquisto della nuova casa.
            </li>
        </ul>
    </div>
</div>

<script lang="ts">
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    export let number: number;
    export let dif: number;
    export let name: string;
</script>

<div
    class="bg-neutral-900 border border-purple-400 p-6 rounded-xl text-center w-full max-w-sm flex-grow"
>
    <h1 class="text-3xl font-extrabold mb-10 leading-tight">{name}</h1>
    <NumberFlowGroup class="flex gap-4 font-extrabold">
        <NumberFlow
            value={number}
            format={{
                style: "currency",
                currency: "EUR",
                maximumFractionDigits: 0,
            }}
            locales={"it-IT"}
            class="text-5xl font-extrabold"
        />
        <br />
        {#if dif}
            <NumberFlow
                value={dif}
                format={{
                    style: "currency",
                    currency: "EUR",
                    signDisplay: "always",
                    maximumFractionDigits: 0,
                }}
                class="text-lg transition-colors duration-300 text-red-500 font-bold"
            />
            <span class="text-xs">vs costo casa</span>
        {/if}
    </NumberFlowGroup>
</div>

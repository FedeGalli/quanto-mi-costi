<script lang="ts">
    import NumberFlow, { NumberFlowGroup } from "@number-flow/svelte";
    export let number: number;
    export let dif: number;
    export let name: string;
    export let sub_text: string;
    export let vs: string;
</script>

<div
    class="bg-neutral-900 border border-purple-400 p-6 rounded-xl text-center w-full max-w-sm flex-grow mb-4"
>
    <h1 class="text-3xl font-extrabold mb-2 leading-tight">{name}</h1>
    <h1 class="text font-extrabold mb-8 leading-tight">{sub_text}</h1>
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
            <span class="text-xs">{vs}</span>
        {/if}
    </NumberFlowGroup>
</div>

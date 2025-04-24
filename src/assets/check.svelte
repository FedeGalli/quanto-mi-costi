<script lang="ts">
    export let label = "";
    export let bind = false;
    export let tooltip = "";
</script>

<div class="checkbox-wrapper-42 flex items-center gap-2">
    <input id={label} type="checkbox" bind:checked={bind} on:change />
    <label class="cbx" for={label}></label>
    <label class="lbl" for={label}><b>{label}</b></label>

    <!-- Info Icon on far right -->
    <div class="relative group ml-auto cursor-pointer">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 192 512"
            class="w-5 h-5 text-gray-500 transform transition-transform duration-200 group-hover:scale-150 group-focus-within:scale-125"
        >
            <path
                d="M48 80a48 48 0 1 1 96 0A48 48 0 1 1 48 80zM0 224c0-17.7 14.3-32 32-32l64 0c17.7 0 32 14.3 32 32l0 224 32 0c17.7 0 32 14.3 32 32s-14.3 32-32 32L32 512c-17.7 0-32-14.3-32-32s14.3-32 32-32l32 0 0-192-32 0c-17.7 0-32-14.3-32-32z"
            />
        </svg>

        <!-- Tooltip -->
        <div
            class="absolute left-full top-1/2 -translate-y-1/2 ml-2 px-3 py-2 text-sm text-white
             bg-gray-700 rounded shadow-lg opacity-0 group-hover:opacity-100
             group-focus-within:opacity-100 transition-opacity duration-300
             pointer-events-none z-10 border border-white max-w"
        >
            {tooltip}
        </div>
    </div>
</div>

<style>
    .checkbox-wrapper-42 input[type="checkbox"] {
        display: none;
        visibility: hidden;
    }

    .checkbox-wrapper-42 {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .checkbox-wrapper-42 .cbx {
        position: relative;
        top: 1px;
        width: 17px;
        height: 17px;
        border: 1px solid #c8ccd4;
        border-radius: 3px;
        vertical-align: middle;
        transition: background 0.1s ease;
        cursor: pointer;
    }
    .checkbox-wrapper-42 .cbx:after {
        content: "";
        position: absolute;
        top: 1px;
        left: 5px;
        width: 5px;
        height: 11px;
        opacity: 0;
        transform: rotate(45deg) scale(0);
        border-right: 2px solid rgb(0, 0, 0);
        border-bottom: 2px solid rgb(0, 0, 0);
        transition: all 0.3s ease;
        transition-delay: 0.15s;
    }
    .checkbox-wrapper-42 .lbl {
        margin-left: 5px;
        vertical-align: middle;
        cursor: pointer;
    }
    .checkbox-wrapper-42 input[type="checkbox"]:checked ~ .cbx {
        border-color: transparent;
        background: #fff;
        animation: jelly-42 0.6s ease;
    }
    .checkbox-wrapper-42 input[type="checkbox"]:checked ~ .cbx:after {
        opacity: 1;
        transform: rotate(45deg) scale(1);
    }
    @-moz-keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
    @-webkit-keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
    @-o-keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
    @keyframes jelly-42 {
        from {
            transform: scale(1, 1);
        }
        30% {
            transform: scale(1.25, 0.75);
        }
        40% {
            transform: scale(0.75, 1.25);
        }
        50% {
            transform: scale(1.15, 0.85);
        }
        65% {
            transform: scale(0.95, 1.05);
        }
        75% {
            transform: scale(1.05, 0.95);
        }
        to {
            transform: scale(1, 1);
        }
    }
</style>

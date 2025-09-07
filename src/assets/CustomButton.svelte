<script lang="ts">
    export let title: string = "";
    export let onClick;
</script>

<button on:click={onClick}>
    <span class="text-base font-semibold leading-tight">{title}</span>
</button>

<style>
    button {
        cursor: pointer;
        border: none;
        padding: 0.9rem 2rem;
        margin: 0;
        font-family: inherit;
        font-size: inherit;
        position: relative;
        display: block; /* Changed from inline-block to block */
        border-radius: 4px;
        overflow: hidden;
        background: rgba(184, 126, 254, 1);
        color: white;
        border: 1px solid transparent;
        width: 100%; /* Add this to ensure full width */
    }

    button span {
        position: relative;
        z-index: 10;
        transition: color 0.4s;
    }

    button:hover span {
        color: rgba(30, 31, 37);
    }

    button::before,
    button::after {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 0;
    }

    button::before {
        content: "";
        background: rgba(20, 21, 26, 1);
        width: 120%;
        left: -10%;
        transform: skew(30deg);
        transition: transform 0.4s cubic-bezier(0.3, 1, 0.8, 1);
    }

    button:hover::before {
        transform: translate3d(100%, 0, 0);
    }
</style>

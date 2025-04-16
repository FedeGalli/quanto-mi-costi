<script lang="ts">
    export let label = "";
    export let bind = false;
</script>

<div class="checkbox-wrapper-42">
    <input id={label} type="checkbox" bind:checked={bind} on:change />
    <label class="cbx" for={label}></label>
    <label class="lbl" for={label}><b>{label}</b></label>
</div>

<style>
    .checkbox-wrapper-42 input[type="checkbox"] {
        display: none;
        visibility: hidden;
    }

    .checkbox-wrapper-42 label {
        display: inline-block;
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

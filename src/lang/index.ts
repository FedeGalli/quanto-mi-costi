import { init, register, locale } from "svelte-i18n";

const defaultLocale = "it";

// Register locales
register("en", () => import("../lang/en.json"));
register("it", () => import("../lang/it.json"));
register("fr", () => import("../lang/fr.json"));

export function switchLanguage(newLocale: string) {
  locale.set(newLocale);
}

export const languages = [
  { code: "en", name: "🇬🇧 English" },
  { code: "it", name: "🇮🇹 Italiano" },
  { code: "fr", name: "🇫🇷 Français" },
];

// Get initial locale from localStorage or browser
function getInitialLocale(): string {
  if (typeof window !== "undefined") {
    // Try to get from localStorage first
    const saved = localStorage.getItem("locale");
    if (saved && ["en", "it", "fr"].includes(saved)) {
      return saved;
    }

    // Fallback to browser language
    const browserLang = navigator.language.split("-")[0];
    return ["en", "it", "fr"].includes(browserLang)
      ? browserLang
      : defaultLocale;
  }
  return defaultLocale;
}

// Initialize i18n
init({
  fallbackLocale: defaultLocale,
  initialLocale: getInitialLocale(),
});

// Save locale changes to localStorage and update HTML lang attribute
if (typeof window !== "undefined") {
  locale.subscribe((value) => {
    if (value) {
      localStorage.setItem("locale", value);
      document.documentElement.lang = value;
    }
  });
}

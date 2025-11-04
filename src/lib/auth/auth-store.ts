// stores/auth.js
import { writable } from "svelte/store";
import { onAuthStateChanged, signOut } from "firebase/auth";
import { auth } from "./credentials";

// Create the user store
export const user = writable(null);
export const isAuthenticated = writable(false);
export const isLoading = writable(true);
export const apiURL =
  //"https://quanto-mi-costi-934184719806.europe-west8.run.app";
  "http://localhost:8080";

export const dataApiURL =
  "https://backend-python-934184719806.europe-west8.run.app";

// Initialize auth state listener
let unsubscribe: any;

export function initAuthStore() {
  if (unsubscribe) return; // Prevent multiple listeners

  unsubscribe = onAuthStateChanged(auth, async (firebaseUser) => {
    if (firebaseUser && firebaseUser.emailVerified) {
      // User is signed in and email is verified
      const userData: any = {
        uid: firebaseUser.uid,
        email: firebaseUser.email,
        displayName: firebaseUser.displayName,
        photoURL: firebaseUser.photoURL,
        emailVerified: firebaseUser.emailVerified,
        firstName: firebaseUser.displayName?.split(" ")[0] || "",
        lastName: firebaseUser.displayName?.split(" ").slice(1).join(" ") || "",
      };

      const pro = await fetch(apiURL + "/is-pro", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          uid: userData.uid,
        }),
      });

      user.set(userData);
      isAuthenticated.set(true);
    } else {
      // User is signed out or email not verified
      user.set(null);
      isAuthenticated.set(false);
    }
    isLoading.set(false);
  });
}

// Logout function
export async function logout() {
  try {
    await signOut(auth);
    user.set(null);
    isAuthenticated.set(false);
  } catch (error) {
    console.error("Logout error:", error);
    throw error;
  }
}

// Cleanup function
export function destroyAuthStore() {
  if (unsubscribe) {
    unsubscribe();
    unsubscribe = null;
  }
}

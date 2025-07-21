// First, install Firebase dependencies:
// npm install firebase

// firebase.js - Firebase configuration
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyCR_qo-F8CQnYEBDU8yssq4Gy4vqdngvwc",
  authDomain: "quanto-mi-costi.firebaseapp.com",
  projectId: "quanto-mi-costi",
  storageBucket: "quanto-mi-costi.firebasestorage.app",
  messagingSenderId: "934184719806",
  appId: "1:934184719806:web:69eb4dfe1d7810b24f15ef",
  measurementId: "G-DTJPQYK5KX",
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);

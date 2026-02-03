import { useEffect, useState } from "react";
import { Routes, Route, Navigate, useNavigate } from "react-router-dom";

import TopBar from "./components/TopBar";
import DecksView from "./components/DecksView";
import AuthModal from "./components/AuthModal";
import { loginUser, registerUser, fetchDecks, createDeck } from "./api";
import DeckDetailView from "./components/DeckDetailView";

function RequireAuth({ token, children }) {
    if (!token) return <Navigate to="/login" replace />;
    return children;
}

export default function App() {
    const navigate = useNavigate();

    const [token, setToken] = useState(localStorage.getItem("token") || "");
    const [authOk, setAuthOk] = useState(false);
    const [decks, setDecks] = useState([]);
    const [msg, setMsg] = useState("");

    // auth inputs
    const [mode, setMode] = useState("login");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("v@v.com");
    const [password, setPassword] = useState("123456");
    const [showPassword, setShowPassword] = useState(false);

    // deck creation inputs
    const [deckTitle, setDeckTitle] = useState("");
    const [deckDescription, setDeckDescription] = useState("");

    // Load decks when token changes
    useEffect(() => {
        let cancelled = false;

        async function run() {
            setAuthOk(false);
            setDecks([]);

            if (!token) return;

            try {
                const data = await fetchDecks(token);

                // normalize shape
                let list = [];
                if (Array.isArray(data)) list = data;
                else if (Array.isArray(data.decks)) list = data.decks;
                else if (Array.isArray(data.data)) list = data.data;

                if (cancelled) return;
                setDecks(list);
                setAuthOk(true);
            } catch (e) {
                if (cancelled) return;
                setAuthOk(false);
                setDecks([]);
                setMsg("Deck load failed: " + (e?.message || String(e)));
            }
        }

        run();
        return () => {
            cancelled = true;
        };
    }, [token]);

    async function login(e) {
        e.preventDefault();
        setMsg("Logging in...");

        try {
            const data = await loginUser(email, password);
            const t = data.token || data.access_token || data.jwt || data.Token;

            if (!t || typeof t !== "string") {
                setMsg("Login succeeded but token field is missing in response.");
                console.log("login response:", data);
                return;
            }

            localStorage.setItem("token", t);
            setToken(t);
            setMsg("");
            navigate("/decks", { replace: true });
        } catch (e) {
            setMsg("Login failed: " + (e?.message || String(e)));
        }
    }

    async function register(e) {
        e.preventDefault();
        setMsg("Registering...");

        try {
            await registerUser(username, email, password);
            setMode("login");
            setMsg("✅ Registered. Now log in.");
        } catch (e) {
            setMsg("Register failed: " + (e?.message || String(e)));
        }
    }

    async function addDeck(e) {
        e.preventDefault();
        if (!token) return;

        const title = deckTitle.trim();
        const description = deckDescription.trim();

        if (!title) {
            setMsg("Deck title is required.");
            return;
        }

        setMsg("Creating deck...");

        try {
            await createDeck(token, title, description);
            setDeckTitle("");
            setDeckDescription("");
            setMsg("✅ Deck created");

            const data = await fetchDecks(token);
            let list = [];
            if (Array.isArray(data)) list = data;
            else if (Array.isArray(data.decks)) list = data.decks;
            else if (Array.isArray(data.data)) list = data.data;

            setDecks(list);
            setAuthOk(true);
        } catch (e) {
            setMsg("Create deck failed: " + (e?.message || String(e)));
        }
    }

    function logout() {
        localStorage.removeItem("token");
        setToken("");
        setAuthOk(false);
        setDecks([]);
        setMsg("Logged out");
        navigate("/login", { replace: true });
    }

    return (
        <div
            style={{
                minHeight: "100vh",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                background: "#0d0d0d",
                color: "#fff",
                padding: 16,
                boxSizing: "border-box",
            }}
        >
            <div
                style={{
                    width: "100%",
                    maxWidth: 900,
                }}
            >
        <div
            style={{
                minHeight: "100vh",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                background: "#0d0d0d",
                color: "#fff",
                padding: 16,
                boxSizing: "border-box",
            }}
        >
            <div
                style={{
                    width: "100%",
                    maxWidth: 900,
                }}
            >
            <TopBar
                token={token}
                onOpenAuth={() => navigate("/login")}
                onLogout={logout}
            />

            <Routes>
                <Route path="/" element={<Navigate to="/decks" replace />} />

                <Route
                    path="/decks"
                    element={
                        <RequireAuth token={token}>
                            <>
                                <DecksView
                                    token={token}
                                    authOk={authOk}
                                    decks={decks}
                                    deckTitle={deckTitle}
                                    setDeckTitle={setDeckTitle}
                                    deckDescription={deckDescription}
                                    setDeckDescription={setDeckDescription}
                                    onCreateDeck={addDeck}
                                />
                                {msg && (
                                    <pre style={{ marginTop: 12, background: "#000000", padding: 12, borderRadius: 8, whiteSpace: "pre-wrap" }}>
                    {msg}
                  </pre>
                                )}
                            </>
                        </RequireAuth>
                    }
                />

                {/* login "page" implemented as modal */}
                <Route
                    path="/login"
                    element={
                        <AuthModal
                            open={true}
                            token={token}
                            mode={mode}
                            setMode={setMode}
                            username={username}
                            setUsername={setUsername}
                            email={email}
                            setEmail={setEmail}
                            password={password}
                            setPassword={setPassword}
                            showPassword={showPassword}
                            setShowPassword={setShowPassword}
                            onLogin={login}
                            onRegister={register}
                            onClose={() => {
                                // If logged in, close goes back to decks, otherwise keep user on login
                                if (token) navigate("/decks");
                            }}
                        />
                    }
                />

                <Route path="*" element={<Navigate to="/decks" replace />} />
                <Route
                    path="/decks/:deckId"
                    element={
                        <RequireAuth token={token}>
                            <DeckDetailView token={token} />
                        </RequireAuth>
                    }
                />
            </Routes>
        </div>
        </div>
            </div>
        </div>
    );
}

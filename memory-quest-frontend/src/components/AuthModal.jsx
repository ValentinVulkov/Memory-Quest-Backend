import React from "react";

export default function AuthModal({
                                      open,
                                      token,
                                      mode,
                                      setMode,
                                      username,
                                      setUsername,
                                      email,
                                      setEmail,
                                      password,
                                      setPassword,
                                      showPassword,
                                      setShowPassword,
                                      onLogin,
                                      onRegister,
                                      onClose,
                                      onGuest,
                                  }) {
    if (!open) return null;

    const input = {
        width: "100%",
        padding: 12,
        boxSizing: "border-box",
        borderRadius: 8,
        border: "1px solid #ccc",
    };

    return (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.5)", display: "flex", justifyContent: "center", alignItems: "center" }}>
            <div style={{ background: "#fff", padding: 16, width: 400, borderRadius: 10 }}>
                {token && <button type="button" onClick={onClose}>Close</button>}

                <form onSubmit={mode === "login" ? onLogin : onRegister} style={{ display: "grid", gap: 10, marginTop: 10 }}>
                    {mode === "register" && (
                        <input value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" style={input} />
                    )}

                    <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" style={input} />

                    <div style={{ display: "flex", gap: 8 }}>
                        <input
                            type={showPassword ? "text" : "password"}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Password"
                            style={{ ...input, flex: 1 }}
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword((s) => !s)}
                            style={{ padding: "0 12px" }}
                        >
                            {showPassword ? "Hide" : "Show"}
                        </button>
                    </div>

                    <button type="submit">{mode === "login" ? "Login" : "Register"}</button>

                    {!token && (
                        <button
                            type="button"
                            onClick={onGuest}
                            style={{ width: "100%", marginTop: 10, padding: 12 }}
                        >
                            Continue as guest
                        </button>
                    )}

                </form>

                <button type="button" onClick={() => setMode(mode === "login" ? "register" : "login")} style={{ marginTop: 10 }}>
                    Switch to {mode === "login" ? "Register" : "Login"}
                </button>
            </div>
        </div>
    );
}

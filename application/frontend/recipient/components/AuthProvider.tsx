'use client';

import { createContext, useContext, useState } from "react";

type Auth = {
    isAuthenticated: false,
} | {
    isAuthenticated: true,
    privateKey: string,
    publicKey: string,
    fullName: string,
};

const AuthContext = createContext<{
    auth: Auth,
    setAuth: (auth: Auth) => void,
}>({
    auth: { isAuthenticated: false },
    setAuth: x => { },
});

export function AuthProvider({ children }: {
    children: React.ReactNode,
}) {
    const [auth, setAuth] = useState<Auth>({
        isAuthenticated: false,
    });

    return <AuthContext.Provider value={{ auth, setAuth }}>
        {children}
    </AuthContext.Provider>;
}

export const useAuth = () => useContext(AuthContext);

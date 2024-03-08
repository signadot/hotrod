import React, {createContext, useContext, ReactNode, useRef} from 'react';

type SessionContextType = {
    sessionID: number | undefined;
    getLastRequestID: () => number;
}
const SessionContext = createContext<SessionContextType>({} as SessionContextType);

interface SessionProviderProps {
    children: ReactNode;
}

const createSessionID = () => {
    return Math.round(Math.random() * 10000);
}

export const SessionProvider: React.FC<SessionProviderProps> = ({ children }) => {
    const sessionIDRef = useRef(createSessionID());
    const lastRequestIDRef = useRef(1);

    const getLastRequestID = () => {
        const current = lastRequestIDRef.current;

        lastRequestIDRef.current += 1;
        return current;
    }

    return (
        <SessionContext.Provider value={{
            sessionID: sessionIDRef.current,
            getLastRequestID
        }}>
            {children}
        </SessionContext.Provider>
    );
};

export const useSession = () => {
    const context = useContext(SessionContext);
    if (context === undefined) {
        throw new Error('useSession must be used within a SessionProvider');
    }
    return context;
};

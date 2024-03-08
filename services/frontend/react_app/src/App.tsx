import './App.css'
import {HomePage} from "./pages/home.tsx";
import {SessionProvider} from "./context/sessionContext/context.tsx";

function App() {
  return (
      <SessionProvider>
        <HomePage />
      </SessionProvider>
  )
}

export default App

import React from 'react';
import BookingPage from "./pages/booking"
import Layout from "./components/layout/layout";
import "./globalStyles.css"
import RouteMap from "./pages/routemap";

function App() {
  return (
    <Layout sidebarContent={<BookingPage />} mainContent={<RouteMap />} />
  );
}

export default App;


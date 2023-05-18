import React from "react";
import styles from "./layout.module.css";
import logo from "../../images/logo.svg";

interface LayoutProps {
  sidebarContent: React.ReactNode;
  mainContent: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({sidebarContent, mainContent}) => {
  return (
    <div className={styles.container}>
      <div className={styles.side}>
        <div className={styles.header}>
          <p className={styles.logo}><img src={logo} alt="" /></p>
          <div className={styles.title}>
            <h2>Bot R.O.D.</h2>
            <span>Ride-On-Demand by Bot drivers</span>
          </div>
        </div>
        <div>
          {sidebarContent}
        </div>
      </div>
      <div className={styles.main}>{mainContent}</div>
    </div>
  );
}

export default Layout;
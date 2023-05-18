import React from "react";
import RouteImage from "../../images/route.png";
import styles from "./routemap.module.css";

const RouteMap = () => (
  <div className={styles.container}>
    <img className={styles.route} src={RouteImage} />
  </div>
)

export default RouteMap;
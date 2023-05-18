import React from "react";
import {RideInfo} from "../../pages/booking/bookingPage";
import styles from "./ride.module.css";

interface RideProps {
  ride: RideInfo;
}

const Ride: React.FC<RideProps> = ({ride}) => {
  return (
    <div className={styles.container}>
      <div>
        <img src={ride.Driver.ImageURL} alt="" />
      </div>
      <div className={styles.info}>
        <p>Your bot driver <span>{ride.Driver.Name}</span> is on the way!</p>
        <p className={styles.eta}>ETA: {ride.ETA/(1000*1000*1000*60)} minutes</p>
      </div>
    </div>
  )
}

export default Ride;
import React from "react";
import Ride from "../../components/ride";
import styles from "./bookingPage.module.css";

export interface Driver {
  DriverID: string;
  Name: string;
  Location: string;
  ImageURL: string;
}

export interface RideInfo {
  Driver: Driver;
  ETA: number;
}

const customerID = 392;

const BookingPage = () => {
  const [rides, setRides] = React.useState<RideInfo[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);

  const from = "15600 NE 8th St, Bellevue, WA, 98008";
  const to = "4400 88th Ave SE, Mercer Island, WA, 98040";
  const addARide = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`/api/dispatch?customer=${customerID}`)
      const ride = await response.json();
      setIsLoading(false);
      setRides([ride, ...rides]);
    } catch (e) {
      setIsLoading(false);
    }
  };
  return (
    <div>
      <div className={styles.tripInfo}>
        <h3 className={styles.title}>Trip Information</h3>
        <div>
          <p><input type="text" className={styles.location} disabled value={from} /></p>
          <p><input type="text" className={styles.location} disabled value={to} /></p>
        </div>
        <div>
          <button type="button" className={styles.book} onClick={addARide} disabled={isLoading}>Add a Ride</button>
        </div>
      </div>
      <div>
        {rides.map((ride, idx) => <Ride key={`ride-${idx}`} ride={ride} /> )}
      </div>
    </div>
  );
}

export default BookingPage;
import React, { useEffect, useState } from "react";
import logo from "./logo.svg";
import "./App.css";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";
import { FactoryTwoTone, SignalWifi2BarTwoTone, WifiChannelTwoTone } from "@mui/icons-material";
import GCILogo from "./assets/gci_logo.png";
import SFSLogo from "./assets/sfs_logo.png";
import SFSNationalLogo from "./assets/sfs_national.png";

interface ReconAP {
  scan_id: number;
  ssid: string;
  bssid: string;
  encryption: number;
  hidden: boolean;
  wps: boolean;
  channel: number;
  signal: number;
  data: number;
  first_seen: number;
  last_seen: number;
  last_seen_delta: number;
  probes: number;
  mfp: boolean;
  clients: {
    ScanID: number;
    client_mac: string;
    ap_mac: string;
    ap_channel: number;
    data: number;
    broadcast_proves: number;
    first_seen: number;
    last_seen: number;
    last_seen_delta: number;
    ssid: string;
  }[];
  num_clients: number;
  vendor: string;
}

enum ReconAPEncryptionType {
  OPEN = 0, // Open
  WEP = 2, // WEP
  WPA2_PSK = 17184063752, // WPA2 (PSK)
  WPA2_MIXED_PSK = 17185112140, // WPA2 Mixed (PSK)
  WPA3_SAE = 2199027450128, // WPA3 (SAE)
}

const MAX_APS = 15;

const prependNewAP = (newAp: ReconAP): ((prevValue: ReconAP[]) => ReconAP[]) => {
  return (prevValue: ReconAP[]): ReconAP[] => {
    if (prevValue.length >= MAX_APS) {
      return [newAp, ...prevValue.slice(0, MAX_APS - 2)];
    }
    return [newAp, ...prevValue];
  };
};

const App: React.FC = (): React.ReactElement => {
  const [apList, setApList] = useState<ReconAP[]>([]);
  useWebSocket(`${process.env.REACT_APP_WS_URL}/api/subscribe`, {
    onMessage: (event) => {
      const newAp = JSON.parse(event.data);
      setApList(prependNewAP(newAp));
    },
  });

  return (
    <div className="App">
      <div className="ap-list">
        {apList.map((ap) => (
          <div key={`ap-${ap.ssid}-${ap.bssid}-${ap.last_seen}`} className="ap">
            <div className="ap-title">
              <span className="ap-ssid">{ap.ssid || "[HIDDEN SSID]"}</span>
              <span className="ap-lastseen">{new Date(ap.last_seen * 1000).toLocaleTimeString()}</span>
            </div>
            <div className="ap-details">
              <span>
                <FactoryTwoTone />
                {ap.vendor}
              </span>
              <span>
                <SignalWifi2BarTwoTone />
                {ap.signal} dBi
              </span>
              <span>
                <WifiChannelTwoTone />
                Channel {ap.channel}
              </span>
            </div>
          </div>
        ))}
      </div>
      <div className="explanation">
        <h1>
          SSID Jungle - <span style={{ color: "#F6EA00" }}>WiFi Pineapple MK7</span>
        </h1>
        <h3>What's this device?</h3>
        <p>
          The device you see on the table in front of you is called a "WiFi Pineapple". No, it definitely doesn't look like a pineapple, but
          against its sweet name it's actually pretty sinister.
        </p>
        <h3>What is it doing?</h3>
        <p>
          The WiFi Pineapple is constantly scanning for WiFi access points (the things you see on the ceiling which you connect to) all
          around it. It is scraping data from these access points and displaying it on screen.
        </p>
        <h3>Why does this matter?</h3>
        <p>
          Using this device, it's possible to learn information about WiFi networks <em>your</em> devices have connected to. With that
          information, this device can then impersonate a network you trust (e.g. MyHouse-5G). Your device could then auto-connect to this
          known WiFi, allowing for further attacks like Person-In-The-Middle, WPA Handshake sniffing/cracking, etc.
        </p>
        <p className="fineprint">
          The WiFi Pineapple is a registered trademark of Hak5 LLC. This project was created for usage at educational or academic events
          only by Bradley Harker. This is a demonstration, no data collected relating to wireless access points will be used with malicious
          intent.
        </p>

        <img className="logo sfs" src={SFSLogo} />
        <img className="logo gci" src={GCILogo} />
      </div>
    </div>
  );
};

export default App;

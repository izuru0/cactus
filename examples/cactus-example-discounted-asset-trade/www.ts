import { BusinessLogicAssetTrade } from "./business-logic-asset-trade";
import { startCactusSocketIOServer } from "@hyperledger/cactus-cmd-socketio-server";
import { initFabricConnector } from "./fabric-connector";
import { initEthereumConnector } from "./ethereum-connector";

async function startBLP() {
  try {
    await initFabricConnector();
    await initEthereumConnector();

    startCactusSocketIOServer({
      id: "guks32pf",
      plugin: new BusinessLogicAssetTrade(),
    });
  } catch (error) {
    console.error("Could not start discounted-asset-trade BLP:", error);
  }
}

startBLP();

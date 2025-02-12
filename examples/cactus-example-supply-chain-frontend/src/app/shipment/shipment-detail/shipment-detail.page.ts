import { v4 as uuidv4 } from "uuid";

import { Component, Inject, Input, OnInit } from "@angular/core";
import {
  UntypedFormGroup,
  UntypedFormBuilder,
  Validators,
} from "@angular/forms";
import { ModalController } from "@ionic/angular";

import { ApiClient } from "@hyperledger/cactus-api-client";

import {
  Shipment,
  Bookshelf,
  DefaultApi as SupplyChainApi,
} from "@hyperledger/cactus-example-supply-chain-business-logic-plugin";

import { Logger, LoggerProvider } from "@hyperledger/cactus-common";
import { QUORUM_DEMO_LEDGER_ID } from "src/constants";

import { AuthConfig } from "../../common/auth-config";

@Component({
  selector: "app-shipment-detail",
  templateUrl: "./shipment-detail.page.html",
  styleUrls: [],
})
export class ShipmentDetailPage implements OnInit {
  private readonly log: Logger;
  private _supplyChainApi: SupplyChainApi | undefined;
  public form: UntypedFormGroup;
  @Input()
  public shipment: Shipment;
  public bookshelves: Bookshelf[];
  public bookshelfIds: string[];

  constructor(
    private readonly baseClient: ApiClient,
    @Inject(QUORUM_DEMO_LEDGER_ID) private readonly quorumLedgerId: string,
    public readonly modalController: ModalController,
    public readonly formBuilder: UntypedFormBuilder,
  ) {
    this.log = LoggerProvider.getOrCreate({ label: "ShipmentDetailPage" });
  }

  public get supplyChainApi(): SupplyChainApi {
    if (!this._supplyChainApi) {
      throw new Error(`InvalidStateError: _supplyChainApi not initialized.`);
    } else {
      return this._supplyChainApi;
    }
  }

  async ngOnInit(): Promise<void> {
    this.log.debug("component initialized.", this.shipment);

    this._supplyChainApi = await this.baseClient.ofLedger(
      this.quorumLedgerId,
      SupplyChainApi,
      {
        baseOptions: {
          headers: { Authorization: `Bearer ${AuthConfig.authToken}` },
        },
      },
    );

    if (!this.shipment) {
      this.shipment = {
        id: uuidv4(),
        bookshelfId: "",
      };
    }
    this.form = this.formBuilder.group({
      id: [this.shipment.id, Validators.required],
      bookshelfId: [this.shipment.bookshelfId, Validators.required],
    });

    await this.loadData();
  }

  private async loadData(): Promise<void> {
    const { data } = await this.supplyChainApi.listBookshelfV1();
    const { data: bookshelves } = data;
    this.bookshelves = bookshelves;
    this.log.debug(`Fetched BambooHarvest data: %o`, bookshelves);
    this.bookshelfIds = this.bookshelves.map((bh) => bh.id);
    this.log.debug(`BambooHarvest IDs: %o`, this.bookshelfIds);
  }

  onClickFormSubmit(value: Shipment): void {
    this.log.debug("form submitted", value);
    this.shipment = value;
    this.modalController.dismiss(this.shipment);
  }

  onClickBtnCancel(): void {
    this.log.debug("form submission cancelled by user");
    this.modalController.dismiss();
  }
}

import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { WalletAssetsService } from '../wallet-assets/wallet-assets.service';

@Controller('wallets/:wallet_id/assets')
export class WalletAssetsController {
  constructor(private readonly walletsService: WalletAssetsService) {}

  @Get()
  all(@Param('wallet_id') wallet_id: string) {
    return this.walletsService.all({ wallet_id });
  }

  @Post()
  create(
    @Param('wallet_id') wallet_id: string,
    @Body() body: { asset_id: string; shares: number },
  ) {
    return this.walletsService.create({ ...body, wallet_id });
  }
}

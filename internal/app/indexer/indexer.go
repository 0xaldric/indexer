package indexer

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
	"github.com/tonindexer/anton/internal/core/repository"
	"github.com/tonindexer/anton/internal/core/repository/account"
	"github.com/tonindexer/anton/internal/core/repository/block"
	"github.com/tonindexer/anton/internal/core/repository/contract"
	"github.com/tonindexer/anton/internal/core/repository/msg"
	"github.com/tonindexer/anton/internal/core/repository/tx"
)

var _ app.IndexerService = (*Service)(nil)

type Service struct {
	*app.IndexerConfig

	blockRepo   core.BlockRepository
	txRepo      core.TransactionRepository
	msgRepo     repository.Message
	accountRepo core.AccountRepository
	contractRepo core.ContractRepository

	run bool
	mx  sync.RWMutex
	wg  sync.WaitGroup
}

func NewService(cfg *app.IndexerConfig) *Service {
	var s = new(Service)

	s.IndexerConfig = cfg

	// validate config
	if s.Workers < 1 {
		s.Workers = 1
	}
	if s.FromBlock < 2 {
		s.FromBlock = 2
	}

	ch, pg := s.DB.CH, s.DB.PG
	s.txRepo = tx.NewRepository(ch, pg)
	s.msgRepo = msg.NewRepository(ch, pg)
	s.blockRepo = block.NewRepository(ch, pg)
	s.accountRepo = account.NewRepository(ch, pg)
	s.contractRepo = contract.NewRepository(pg)

	return s
}

func (s *Service) running() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.run
}

func (s *Service) Start() error {
	ctx := context.Background()

	fromBlock := s.FromBlock

	lastMaster, err := s.blockRepo.GetLastMasterBlock(ctx)
	switch {
	case err == nil:
		fromBlock = lastMaster.SeqNo + 1
	case err != nil && !errors.Is(err, core.ErrNotFound):
		return errors.Wrap(err, "cannot get last masterchain block")
	}

	s.mx.Lock()
	s.run = true
	s.mx.Unlock()

	blocksChan := make(chan *core.Block, s.Workers*2)
	msgChan := make(chan *core.Message, s.Workers*100000)

	s.wg.Add(1)
	go s.fetchMasterLoop(fromBlock, blocksChan)

	s.wg.Add(1)
	go s.saveBlocksLoop(blocksChan, msgChan)

	s.wg.Add(1)
	go s.produceMessageLoop(msgChan)

	log.Info().
		Uint32("from_block", fromBlock).
		Int("workers", s.Workers).
		Msg("started indexer service")

	return nil
}

func (s *Service) Stop() {
	s.mx.Lock()
	s.run = false
	s.mx.Unlock()

	s.wg.Wait()
}

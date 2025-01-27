package testfactory

import (
	crawlJobRepo "crawlquery/api/crawl/job/repository/mem"
	crawlLogRepo "crawlquery/api/crawl/log/repository/mem"
	crawlService "crawlquery/api/crawl/service"
	crawlThrottleService "crawlquery/api/crawl/throttle/service"
	"crawlquery/api/domain"
	"time"

	nodeRepo "crawlquery/api/node/repository/mem"
	nodeService "crawlquery/api/node/service"

	linkRepo "crawlquery/api/link/repository/mem"
	linkService "crawlquery/api/link/service"

	pageRepo "crawlquery/api/page/repository/mem"
	pageService "crawlquery/api/page/service"

	shardRepo "crawlquery/api/shard/repository/mem"
	shardService "crawlquery/api/shard/service"

	indexJobRepo "crawlquery/api/index/job/repository/mem"
	indexLogRepo "crawlquery/api/index/log/repository/mem"
	indexService "crawlquery/api/index/service"

	eventService "crawlquery/api/event/service"

	"crawlquery/pkg/testutil"
)

type ServiceFactory struct {
	EventService         *eventService.Service
	ShardRepo            *shardRepo.Repository
	ShardService         *shardService.Service
	PageRepo             *pageRepo.Repository
	PageService          *pageService.Service
	LinkRepo             *linkRepo.Repository
	LinkService          *linkService.Service
	NodeRepo             *nodeRepo.Repository
	NodeService          *nodeService.Service
	CrawlJobRepo         *crawlJobRepo.Repository
	CrawlLogRepo         *crawlLogRepo.Repository
	CrawlThrottleService *crawlThrottleService.Service
	CrawlService         *crawlService.Service
	IndexLogRepo         *indexLogRepo.Repository
	IndexJobRepo         *indexJobRepo.Repository
	IndexService         *indexService.Service
}

type ServiceFactoryOption func(*ServiceFactory)

func WithShard(shard *domain.Shard) ServiceFactoryOption {
	return func(factory *ServiceFactory) {
		factory.ShardRepo.Create(shard)
	}
}

func WithNode(node *domain.Node) ServiceFactoryOption {
	return func(factory *ServiceFactory) {
		factory.NodeRepo.Create(node)
	}
}

func WithCrawlJob(job *domain.CrawlJob) ServiceFactoryOption {
	return func(factory *ServiceFactory) {
		factory.CrawlJobRepo.Save(job)
	}
}

func NewServiceFactory(options ...ServiceFactoryOption) *ServiceFactory {

	eventService := eventService.NewService()

	shardRepo := shardRepo.NewRepository()
	shardService := shardService.NewService(
		shardService.WithRepo(shardRepo),
		shardService.WithLogger(testutil.NewTestLogger()),
	)

	nodeRepo := nodeRepo.NewRepository()
	nodeService := nodeService.NewService(
		nodeService.WithNodeRepo(nodeRepo),
		nodeService.WithLogger(testutil.NewTestLogger()),
		nodeService.WithShardService(shardService),
		nodeService.WithRandSeed(time.Now().UnixNano()),
	)

	linkRepo := linkRepo.NewRepository()
	linkService := linkService.NewService(
		linkService.WithEventService(eventService),
		linkService.WithLinkRepo(linkRepo),
		linkService.WithLogger(testutil.NewTestLogger()),
	)

	crawlRepo := crawlJobRepo.NewRepository()
	crawlLogRepo := crawlLogRepo.NewRepository()
	crawlThrottleService := crawlThrottleService.NewService(
		crawlThrottleService.WithRateLimit(time.Second * 20),
	)
	crawlService := crawlService.NewService(
		crawlService.WithEventService(eventService),
		crawlService.WithCrawlThrottleService(crawlThrottleService),
		crawlService.WithCrawlJobRepo(crawlRepo),
		crawlService.WithNodeService(nodeService),
		crawlService.WithCrawlLogRepo(
			crawlLogRepo,
		),
		crawlService.WithLogger(testutil.NewTestLogger()),
		crawlService.WithWorkers(10),
		crawlService.WithMaxQueueSize(100),
	)

	pageRepo := pageRepo.NewRepository()
	pageService := pageService.NewService(
		pageService.WithEventService(eventService),
		pageService.WithPageRepo(pageRepo),
		pageService.WithLogger(testutil.NewTestLogger()),
		pageService.WithShardService(shardService),
	)

	indexJobRepo := indexJobRepo.NewRepository()
	indexLogRepo := indexLogRepo.NewRepository()
	indexService := indexService.NewService(
		indexService.WithEventService(eventService),
		indexService.WithIndexJobRepo(indexJobRepo),
		indexService.WithIndexLogRepo(indexLogRepo),
		indexService.WithLogger(testutil.NewTestLogger()),
	)

	factory := &ServiceFactory{
		EventService:         eventService,
		ShardRepo:            shardRepo,
		ShardService:         shardService,
		PageRepo:             pageRepo,
		PageService:          pageService,
		LinkRepo:             linkRepo,
		LinkService:          linkService,
		NodeRepo:             nodeRepo,
		NodeService:          nodeService,
		CrawlJobRepo:         crawlRepo,
		CrawlLogRepo:         crawlLogRepo,
		CrawlThrottleService: crawlThrottleService,
		CrawlService:         crawlService,
		IndexLogRepo:         indexLogRepo,
		IndexJobRepo:         indexJobRepo,
		IndexService:         indexService,
	}
	for _, option := range options {
		option(factory)
	}
	return factory
}

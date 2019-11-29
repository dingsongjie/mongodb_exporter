package mongod

import (
	"context"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileStatus struct {
	Name  string `bson:"database"`
	Count int64  `bson:"count"`
}

var (
	profileCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "profile",
		Name:      "slow_query_count",
		Help:      "The number of slow queries in this database",
	}, []string{"database"})
)

func (profileStatus *ProfileStatus) Export(ch chan<- prometheus.Metric) {
	profileCount.WithLabelValues(profileStatus.Name).Set(float64(profileStatus.Count))
	profileCount.Collect(ch)
}

func CollectProfileStatus(client *mongo.Client, ch chan<- prometheus.Metric) {

	datebaseNames, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		log.Errorf("Failed to get database names, %v", err)
		return
	}
	for _, db := range datebaseNames {
		if db == "admin" || db == "local" || db == "config" {
			continue
		}
		count, err := client.Database(db).Collection("system.profile").CountDocuments(context.TODO(), bson.M{})
		if err != nil {
			glog.Error(err)
			return
		}
		profileStatus := ProfileStatus{db, count}
		profileStatus.Export(ch)
	}
}

// Describe describes the metrics for prometheus
func (ProfileStatus *ProfileStatus) Describe(ch chan<- *prometheus.Desc) {
	ProfileStatus.Describe(ch)
}

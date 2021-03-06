//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2020 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package kinds

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/usecases/traverser"
)

// DeleteAction Class Instance from the conncected DB
func (m *Manager) DeleteAction(ctx context.Context, principal *models.Principal, id strfmt.UUID) error {
	err := m.authorizer.Authorize(principal, "delete", fmt.Sprintf("actions/%s", id.String()))
	if err != nil {
		return err
	}

	unlock, err := m.locks.LockConnector()
	if err != nil {
		return NewErrInternal("could not acquire lock: %v", err)
	}
	defer unlock()

	return m.deleteActionFromRepo(ctx, id)
}

func (m *Manager) deleteActionFromRepo(ctx context.Context, id strfmt.UUID) error {
	actionRes, err := m.getActionFromRepo(ctx, id, traverser.UnderscoreProperties{})
	if err != nil {
		return err
	}

	action := actionRes.Action()
	err = m.vectorRepo.DeleteAction(ctx, action.Class, id)
	if err != nil {
		return NewErrInternal("could not delete action from vector repo: %v", err)
	}

	return nil
}

// DeleteThing Class Instance from the conncected DB
func (m *Manager) DeleteThing(ctx context.Context, principal *models.Principal, id strfmt.UUID) error {
	err := m.authorizer.Authorize(principal, "delete", fmt.Sprintf("things/%s", id.String()))
	if err != nil {
		return err
	}

	unlock, err := m.locks.LockConnector()
	if err != nil {
		return NewErrInternal("could not acquire lock: %v", err)
	}
	defer unlock()

	return m.deleteThingFromRepo(ctx, id)
}

func (m *Manager) deleteThingFromRepo(ctx context.Context, id strfmt.UUID) error {
	thingRes, err := m.getThingFromRepo(ctx, id, traverser.UnderscoreProperties{})
	if err != nil {
		return err
	}

	thing := thingRes.Thing()
	err = m.vectorRepo.DeleteThing(ctx, thing.Class, id)
	if err != nil {
		return NewErrInternal("could not delete thing from vector repo: %v", err)
	}

	return nil
}

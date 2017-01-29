package main

import (
	"github.com/leeola/kala/client"
	"github.com/leeola/kala/contenttype/inventory"
	"github.com/leeola/kala/index"
	"github.com/leeola/kala/store"
	"github.com/urfave/cli"
)

func listCommand(ctx *cli.Context) error {
	nameOrAnchor := ctx.Args().Get(0)

	if nameOrAnchor == "" {
		return cli.ShowCommandHelp(ctx, "list")
	}

	c, err := ClientFromContext(ctx)
	if err != nil {
		return err
	}

	results, err := queryNameOrAnchor(c, nameOrAnchor)
	if err != nil {
		return err
	}

	if len(results.Hashes) == 0 {
		Printlnf("No results found for %q", nameOrAnchor)
	}

	for _, h := range results.Hashes {
		var v store.Version
		if err := c.GetBlobAndUnmarshal(h.Hash, &v); err != nil {
			return err
		}

		if err := listForVersion(c, v); err != nil {
			return err
		}
	}

	return nil
}

func queryNameOrAnchor(c *client.Client, na string) (index.Results, error) {
	// first, try a name query, as that is more likely to be used.
	q := index.Query{Metadata: index.Metadata{
		"contentType": "inventory",
		"name":        na,
	}}
	results, err := c.Query(q)
	if err != nil {
		return index.Results{}, err
	}

	// if we got results, return them.
	if len(results.Hashes) > 0 {
		return results, nil
	}

	// we got no name matches, so try by anchor
	q = index.Query{Metadata: index.Metadata{
		"contentType": "inventory",
		"container":   `"` + na + `"`,
	}}
	return c.Query(q)
}

func listForVersion(c *client.Client, v store.Version) error {
	var vMeta inventory.Meta
	if err := c.GetBlobAndUnmarshal(v.Meta, &vMeta); err != nil {
		return err
	}

	Printlnf("%s items:", vMeta.Name)

	q := index.Query{Metadata: index.Metadata{
		"contentType": "inventory",
		"container":   `"` + v.Anchor + `"`,
	}}
	results, err := c.Query(q)
	if err != nil {
		return err
	}

	for _, h := range results.Hashes {
		var v store.Version
		if err := c.GetBlobAndUnmarshal(h.Hash, &v); err != nil {
			return err
		}

		var m inventory.Meta
		if err := c.GetBlobAndUnmarshal(v.Meta, &m); err != nil {
			return err
		}

		Printlnf("  %s", m.Name)
	}

	return nil
}

package zinc

import (
	"database/sql"
)

type Args []any
type NamedArg = sql.NamedArg
type NamedArgs map[string]any
type UnitedArgs struct {
	Unnamed Args
	Named   NamedArgs
}

func Named(name string, value any) NamedArg {
	return NamedArg{Name: name, Value: value}
}

func United(args ...any) (UnitedArgs, OptionsModifier) {
	var target UnitedArgs
	var optionsModifiers []OptionsModifier
	ensureNamedArgs := func() {
		if target.Named == nil {
			target.Named = NamedArgs{}
		}
	}
	for _, arg := range args {
		switch a := arg.(type) {
		case NamedArg:
			ensureNamedArgs()
			if a.Name != "" {
				target.Named[a.Name] = a.Value
			} else {
				target.Unnamed = append(target.Unnamed, a.Value)
			}
		case *NamedArg:
			if a != nil {
				ensureNamedArgs()
				if a.Name != "" {
					target.Named[a.Name] = a.Value
				} else {
					target.Unnamed = append(target.Unnamed, a.Value)
				}
			}
		case NamedArgs:
			ensureNamedArgs()
			for k, v := range a {
				target.Named[k] = v
			}
		case map[string]any:
			ensureNamedArgs()
			for k, v := range a {
				target.Named[k] = v
			}
		case UnitedArgs:
			if a.HasUnnamed() {
				target.Unnamed = append(target.Unnamed, a.Unnamed...)
			}
			if a.HasNamed() {
				ensureNamedArgs()
				for k, v := range a.Named {
					target.Named[k] = v
				}
			}
		case *UnitedArgs:
			if a != nil {
				if a.HasUnnamed() {
					target.Unnamed = append(target.Unnamed, a.Unnamed...)
				}
				if a.HasNamed() {
					ensureNamedArgs()
					for k, v := range a.Named {
						target.Named[k] = v
					}
				}
			}
		case OptionsModifier:
			if a != nil {
				optionsModifiers = append(optionsModifiers, a)
			}
		case func(*Options):
			if a != nil {
				optionsModifiers = append(optionsModifiers, a)
			}
		default:
			target.Unnamed = append(target.Unnamed, a)
		}
	}
	if len(optionsModifiers) <= 0 {
		return target, nil
	}
	return target, func(opts *Options) {
		for _, modifier := range optionsModifiers {
			modifier(opts)
		}
	}
}

func (uArgs UnitedArgs) Empty() bool {
	return len(uArgs.Unnamed) <= 0 && len(uArgs.Named) <= 0
}

func (uArgs UnitedArgs) HasUnnamed() bool {
	return len(uArgs.Unnamed) > 0
}

func (uArgs UnitedArgs) HasNamed() bool {
	return len(uArgs.Named) > 0
}
